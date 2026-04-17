"use client";

import { ScoreGraph } from "@/components/charts/ScoreGraph";
import { ScoreGraphSkeleton } from "@/components/skeletons/scoreboard";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { PageNavigation } from "@/components/ui/page-navigation";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { ChartData } from "@/models/score";
import { useEventStore, useScoreStore } from "@/store";
import { toChartData } from "@/utils/scoreTransform";
import { Loader2, Trophy } from "lucide-react";
import { useEffect, useState } from "react";

export default function Scoreboard() {
  const [searchQuery, setSearchQuery] = useState("");
  const [chartData, setChartData] = useState<ChartData>({ labels: [], points: [] });

  const {
    graphLoading,
    scoresLoading,
    currentPage,
    totalPages,
    scores,
    graphScores,
    fetchPage,
    fetchGraph,
  } = useScoreStore();
  const {
    info: { event_start },
  } = useEventStore();

  useEffect(() => {
    fetchPage(currentPage);
    fetchGraph();
  }, []);

  useEffect(() => {
    setChartData(toChartData(graphScores));
  }, [graphScores, event_start]);

  const handlePageChange = async (newPage: number) => {
    if (newPage < 1 || newPage > totalPages) return;
    await fetchPage(newPage);
  };

  const filteredScores = scores.filter((score) =>
    score.team_name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="container p-4 space-y-4">
      {graphLoading ? (
        <ScoreGraphSkeleton />
      ) : (
        graphScores.length !== 0 && <ScoreGraph data={chartData} />
      )}

      <div className="flex flex-col sm:flex-row justify-between items-center gap-4">
        <Input
          type="text"
          placeholder="Search teams..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="max-w-sm"
        />
        <PageNavigation
          currentPage={currentPage}
          totalPages={totalPages}
          onPageChange={handlePageChange}
        />
      </div>

      <Card>
        <CardContent className="p-0">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-24 text-center">Rank</TableHead>
                <TableHead className="text-center">Team Name</TableHead>
                <TableHead className="w-24 text-center">Score</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {scoresLoading ? (
                <TableRow>
                  <TableCell colSpan={3} className="h-32 text-center">
                    <Loader2 className="h-6 w-6 animate-spin mx-auto text-muted-foreground" />
                  </TableCell>
                </TableRow>
              ) : filteredScores.length === 0 ? (
                <TableRow>
                  <TableCell colSpan={3} className="h-32 text-center text-muted-foreground">
                    {searchQuery ? "No teams match your search." : "No scores yet."}
                  </TableCell>
                </TableRow>
              ) : (
                filteredScores.map((entry) => (
                  <TableRow key={entry.team_id}>
                    <TableCell className="text-center">
                      <div className="flex justify-center items-center">
                        {entry.rank <= 3 ? (
                          <Trophy
                            className={`w-6 h-6 ${
                              entry.rank === 1
                                ? "text-yellow-500"
                                : entry.rank === 2
                                  ? "text-gray-400"
                                  : "text-orange-500"
                            }`}
                          />
                        ) : (
                          <Badge variant="secondary" className="w-8 flex justify-center">
                            #{entry.rank}
                          </Badge>
                        )}
                      </div>
                    </TableCell>
                    <TableCell className="font-medium text-center">
                      <span className="truncate block max-w-xs mx-auto">{entry.team_name}</span>
                    </TableCell>
                    <TableCell className="text-center">{entry.score}</TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </CardContent>
      </Card>

      <PageNavigation
        currentPage={currentPage}
        totalPages={totalPages}
        onPageChange={handlePageChange}
      />
    </div>
  );
}
