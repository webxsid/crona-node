import React from "react";
import { Box, Text } from "ink";
import type { ActiveContext } from "../interfaces/context.interface.js";

export function Header({ context }: { context?: ActiveContext | null }) {
  return (
    <Box borderStyle="single" paddingX={1} >
      <Text>
        Repo: {context?.repoId ?? "-"} | Stream: {context?.streamId ?? "-"} | Issue: {context?.issueId ?? "-"}
      </Text>
    </Box>
  );
}
