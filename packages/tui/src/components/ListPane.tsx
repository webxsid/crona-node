import React from "react";
import { Box, Text } from "ink";
import { Pane } from "./Pane.js";

interface Props<T> {
  title: string;
  items: T[];
  cursor: number;
  paneActive: boolean;
  renderItem: (item: T) => string;
}

export function ListPane<T>({
  title,
  items,
  cursor,
  paneActive,
  renderItem
}: Props<T>) {

  return (
    <Pane title={title} active={paneActive}>
      <Box flexDirection="column" paddingX={1}>
        {items.map((item, i) => (
          <Text
            key={i}
            color={paneActive && i === cursor ? "green" : undefined}
          >
            {paneActive && i === cursor ? "▶ " : "  "}
            {renderItem(item)}
          </Text>
        ))}
      </Box>
    </Pane>
  );
}
