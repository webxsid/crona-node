export function clearTerminal() {
  process.stdout.write("\x1Bc");
}
