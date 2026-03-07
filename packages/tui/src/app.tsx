import React, { useEffect, useState } from "react";
import { Box, Text, useApp, useInput, useStdout } from "ink";
import { StatusView } from "./views/status.js";
import { useKernelEvents } from "./hooks/useKernelEvents.js";
import { useTuiState } from "./hooks/useTuiState.js";
import { useKeybindings } from "./hooks/useKeybindings.js";
import { useRepos } from "./hooks/useRepo.js";
import { useStreams } from "./hooks/useStream.js";
import { useIssues } from "./hooks/useIssue.js";
import { useScratchpads } from "./hooks/useScratchPads.js";
import { useContext } from "./hooks/useContext.js";
import { Layout } from "./layout/index.js";
import { Header } from "./layout/Header.js";
import { ListPane } from "./components/ListPane.js";
import { LeftPanel } from "./layout/LeftPanel.js";
import { HelpBar } from "./components/HelpBar.js";

export interface AppProps {
  kernel: {
    baseUrl: string;
    token: string;
  };
}

export function App({ kernel }: AppProps) {

  useKernelEvents(kernel);

  const tui = useTuiState();

  const context = useContext(kernel);

  const repos = useRepos(kernel)
  const streams = useStreams(kernel, context?.repoId);
  const issues = useIssues(kernel, context?.streamId);
  const scratchpads = useScratchpads(kernel);

  useKeybindings(
    kernel,
    {
      pane: tui.state.pane,
      moveUp: tui.moveUp,
      moveDown: tui.moveDown,
      focusNextPane: tui.focusNextPane,
      focusPrevPane: tui.focusPrevPane,
      focusPane: tui.focusPane,
      setCommandBuffer: tui.setCommandBuffer,
      clearCommand: tui.clearCommand
    },
    () => {
      if (tui.state.pane === "repos") return repos.length;
      if (tui.state.pane === "streams") return streams.length;
      if (tui.state.pane === "issues") return issues.length;
      if (tui.state.pane === "scratchpads") return scratchpads.length;
      return 0;
    }
  );

  const { exit } = useApp();
  const { stdout } = useStdout()

  const [size, setSize] = useState<{
    width: number;
    height: number;
  }>({
    width: stdout.columns ?? 80,
    height: stdout.rows ?? 24,
  });

  const { width, height } = size;

  const clearAndExit = () => {
    exit();
  }

  // Handle terminal resize
  useEffect(() => {
    function onResize() {
      setSize({
        width: stdout.columns ?? 80,
        height: stdout.rows ?? 24,
      });
    }

    stdout.on("resize", onResize);
    return () => {
      stdout.off("resize", onResize);
    };
  }, [stdout]);

  useInput((input) => {
    if (input === "q") clearAndExit();
  });


  useEffect(() => {
    console.log("Crona TUI started. Press 'q' to quit.");
    return () => {
      console.log("Crona TUI exited.");
    }
  }, [])

  return (
    <Layout>

      <Header context={context} />

      <Box flexGrow={1}>

        {/* LEFT SIDE */}
        <LeftPanel width={35}>

          {/* Repo + Stream row */}
          <Box flexDirection="row" height="30%">
            <ListPane
              title="Repos [1]"
              items={repos}
              cursor={tui.state.cursor.repos}
              paneActive={tui.state.pane === "repos"}
              renderItem={(r) => r.name}
            />

            <ListPane
              title="Streams [2]"
              items={streams}
              cursor={tui.state.cursor.streams}
              paneActive={tui.state.pane === "streams"}
              renderItem={(s) => s.name}
            />
          </Box>

          {/* Issues below */}
          <Box flexGrow={1}>
            <ListPane
              title="Issues [3]"
              items={issues}
              cursor={tui.state.cursor.issues}
              paneActive={tui.state.pane === "issues"}
              renderItem={(i) => i.title}
            />
          </Box>

        </LeftPanel>


        {/* RIGHT SIDE */}
        <Box flexDirection="column" flexGrow={1}>
          <ListPane
            title="Scratchpads [4]"
            items={scratchpads}
            cursor={tui.state.cursor.scratchpads}
            paneActive={tui.state.pane === "scratchpads"}
            renderItem={(s) => s.name}
          />
        </Box>

      </Box>

      <HelpBar pane={tui.state.pane} />
    </Layout>
  );
}
