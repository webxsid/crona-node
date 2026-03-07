import React from "react";
import { Box } from "ink";

export function LeftPanel({ children, width }: { children: React.ReactNode; width: number }) {
  return (
    <Box width={`${width}%`} flexDirection="column">
      {children}
    </Box>
  );
}
