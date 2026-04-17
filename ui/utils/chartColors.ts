/**
 * Generate an array of evenly-spaced HSL colors for chart series.
 */
export function generateChartColors(count: number): string[] {
  if (count === 0) return [];
  const step = 360 / count;
  return Array.from({ length: count }, (_, i) => `hsl(${i * step}, 70%, 50%)`);
}

/**
 * Get the HSL color for a single chart series by index.
 */
export function chartColor(index: number, total: number): string {
  const step = 360 / Math.max(total, 1);
  return `hsl(${index * step}, 70%, 50%)`;
}
