import * as program from "commander";

import downloadPlex from "./downloadPlex";

program
  .command("download <plexClaim> <releaseChannel> <outDir>")
  .description("download Plex Media Server")
  .action(downloadPlex);

program.parse(process.argv);
