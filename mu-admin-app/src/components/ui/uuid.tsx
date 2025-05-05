import React from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./tooltip";
import { cn } from "@/lib/utils";

type Props = {
  uuid: string;
  length?: number;
  className?: string;
};

export const UUID = ({ uuid, length = 8, className }: Props) => {
  return (
    <TooltipProvider>
      <Tooltip>
        <TooltipTrigger>
          <span className={cn(className)}>
            {uuid.substring(0, length)}...{uuid.substring(uuid.length - 3)}
          </span>
        </TooltipTrigger>
        <TooltipContent>
          <p>{uuid}</p>
        </TooltipContent>
      </Tooltip>
    </TooltipProvider>
  );
};
