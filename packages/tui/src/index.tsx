#!/usr/bin/env node

import React from "react";
import { render } from "ink";
import { Bootstrap } from "./boorstrap.js";

const { unmount } = render(
  <Bootstrap unmount={unmountHandler} />,
  {
    stdin: process.stdin,
    stdout: process.stdout,
    exitOnCtrlC: true,
  }
);

function unmountHandler() {
  unmount();
}

const clearAndExit = () => {
  unmount();
  process.stdout.write("\x1Bc");
  process.exit(0);
};

process.on("SIGINT", clearAndExit);
process.on("SIGTERM", clearAndExit);
process.on("exit", () => {
  process.stdout.write("\x1Bc");
});
