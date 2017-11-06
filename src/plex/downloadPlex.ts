import * as fs from "fs";
import * as path from "path";
import fetch from "node-fetch";
import * as ProgressBar from "progress";

import exitError from "./exitError";

type Release = "plexpass" | "public";

export default async (plexClaim: string, release: Release, outDir: string) => {
  if (!fs.existsSync(outDir)) exitError(`Directory ${outDir} does not exist.`);

  if (release !== "plexpass" && release !== "public")
    exitError("Release channel must be plexpass or public.");

  const releaseInfo = await fetchReleaseInfo(plexClaim, release, outDir);

  if (fs.existsSync(releaseInfo.filePath)) return; // TODO: Check hash
  await fetchRelease(releaseInfo);
};

const fetchRelease = async (releaseInfo: ReleaseInfo) => {
  const file = fs.createWriteStream(releaseInfo.filePath);
  const res = await fetch(releaseInfo.fileUrl);
  const contentLength = parseInt(res.headers.get("Content-Length"), 10);
  const bar = new ProgressBar(
    `Downloading Plex ${releaseInfo.version} [:bar]`,
    {
      complete: "=",
      incomplete: " ",
      width: 20,
      total: contentLength
    }
  );
  res.body.on("data", chunk => {
    bar.tick(chunk.length);
    file.write(chunk);
  });
  res.body.on("end", () => {
    console.log("\n");
    file.close();
  });
};

interface ReleaseInfo {
  version: string;
  fileHash: string;
  fileName: string;
  fileUrl: string;
  filePath: string;
}

const fetchReleaseInfo = async (
  plexClaim: string,
  release: Release,
  outDir: string
) => {
  const releaseCheckBaseUrl =
    "https://plex.tv/downloads/details/1?build=linux-ubuntu-x86_64";
  let releaseCheckUrl = releaseCheckBaseUrl;
  releaseCheckUrl += `&channel=${release === "plexpass" ? "8" : "16"}`;
  releaseCheckUrl += "&distro=ubuntu";
  if (release === "plexpass") releaseCheckUrl += "&X-Plex-Token=" + plexClaim;

  const releaseCheckResponse = await fetch(releaseCheckUrl);
  const releaseCheckXml = await releaseCheckResponse.text();
  const fileName = parseFileName(releaseCheckXml);
  const releaseInfo: ReleaseInfo = {
    fileName,
    version: parseVersion(releaseCheckXml),
    fileHash: parseFileHash(releaseCheckXml),
    fileUrl: `https://plex.tv${parseFile(releaseCheckXml)}`,
    filePath: path.resolve(outDir, fileName)
  };

  return releaseInfo;
};

const parseFromXml = (regex: RegExp, error: string, xml: string) => {
  const matches = xml.match(regex);
  if (matches === null || matches.length < 2) exitError(error);
  return matches![1];
};

type XmlParser = (releaseXml: string) => string;

const parseVersion: XmlParser = parseFromXml.bind(
  null,
  /.*Release.*version="([^"]*)"/,
  "Unable to retrieve Plex release version."
);

const parseFileHash: XmlParser = parseFromXml.bind(
  null,
  /.*fileHash="([^"]*)"/,
  "Unable to retrieve release file hash."
);

const parseFileName: XmlParser = parseFromXml.bind(
  null,
  /.*fileName="([^"]*)"/,
  "Unable to retrieve release filename."
);

const parseFile: XmlParser = parseFromXml.bind(
  null,
  /.*file="([^"]*)"/,
  "Unable to retrieve release url."
);
