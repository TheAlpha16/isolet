import { Button } from "@/components/ui/button";
import { ChevronLeft, ChevronRight } from "lucide-react";

interface PageNavigationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

export function PageNavigation({ currentPage, totalPages, onPageChange }: PageNavigationProps) {
  return (
    <div className="flex items-center justify-center gap-1">
      <Button
        onClick={() => onPageChange(currentPage - 1)}
        disabled={currentPage === 1}
        variant="ghost"
        size="icon"
        className="rounded-full"
      >
        <ChevronLeft className="w-5 h-5" />
      </Button>
      <span className="text-sm text-muted-foreground px-2">
        Page {currentPage} of {totalPages}
      </span>
      <Button
        onClick={() => onPageChange(currentPage + 1)}
        disabled={currentPage === totalPages}
        variant="ghost"
        size="icon"
        className="rounded-full"
      >
        <ChevronRight className="w-5 h-5" />
      </Button>
    </div>
  );
}
