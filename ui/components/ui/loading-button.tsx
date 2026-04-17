import { Button } from "@/components/ui/button";
import { Loader2 } from "lucide-react";

type LoadingButtonProps = React.ComponentProps<typeof Button> & {
  loading?: boolean;
};

export function LoadingButton({ loading, children, disabled, ...props }: LoadingButtonProps) {
  return (
    <Button disabled={loading || disabled} {...props}>
      {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
      {children}
    </Button>
  );
}
