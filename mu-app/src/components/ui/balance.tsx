import { cn } from "@/lib/utils";

type Props = {
  balance: number;
};

export const Balance = ({ balance }: Props) => {
  return (
    <p
      className={cn(
        "text-lg font-medium",
        balance < 0 ? "text-red-500" : "text-green-500"
      )}
    >
      {balance} ₽
    </p>
  );
};
