// A single point in the chart
export interface ChartPoint {
  timestamp: number; // X-axis
  [team_name: string]: number; // teamName -> score
}

// Labels extracted alongside points
export interface ChartData {
  labels: string[]; // e.g. ["Team A", "Team B"]
  points: ChartPoint[];
}

export interface CategoryProgress {
  category: string;
  solved: number;
  total: number;
}
