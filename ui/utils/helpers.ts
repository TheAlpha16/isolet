export function fromUnixSeconds(sec: number): Date {
  return new Date(sec * 1000);
}
