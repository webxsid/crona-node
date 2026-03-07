import React from "react";
import { Box, Text } from "ink";

interface PaneProps {
  title: string
  active?: boolean
  children?: React.ReactNode
}

export function Pane({ title, active, children }: PaneProps) {
  return (
    <Box
      flexDirection="column"
      borderStyle="single"
      borderColor={active ? "green" : "gray"}
      flexGrow={1}
      position="relative"
    >
      <Box marginLeft={1} marginTop={-1} position="absolute" backgroundColor={"black"} paddingX={1}>
        <Text color={active ? "green" : "gray"}>
          {title}
        </Text>
      </Box>

      <Box flexDirection="column" flexGrow={1}>
        {children}
      </Box>
    </Box>
  );
}
