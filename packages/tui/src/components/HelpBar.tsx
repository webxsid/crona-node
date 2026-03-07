import React from "react";
import { Box, Text } from "ink";
import { getPaneActions } from "../config/pane-registry.js";
import type { Pane } from "../config/pane-registry.js";

interface Props {
  pane: Pane;
}

export function HelpBar({ pane }: Props) {
  const actions = getPaneActions(pane);

  if (!actions.length) return null;

  return (
    <Box paddingX={1} justifyContent="space-between">

      {/* LEFT: pane actions */}
      <Box>
        {actions.map((action, i) => (
          <Text key={action.id}>
            <Text color="cyan">[{action.key}]</Text> {action.label}
            {i !== actions.length - 1 ? "   " : ""}
          </Text>
        ))}
      </Box>

      {/* RIGHT: global actions */}
      <Box>
        <Text>
          <Text color="cyan">[q]</Text> quit
        </Text>
      </Box>

    </Box>
  );
}
