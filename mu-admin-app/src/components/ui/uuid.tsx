import React from "react";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "./tooltip";

type Props = {
  uuid: string;
  length?: number;
};

export const UUID = ({ uuid, length = 8 }: Props) => {
  return (
    <div>
      <TooltipProvider>
        <Tooltip>
          <TooltipTrigger>
            {uuid.substring(0, length)}...{uuid.substring(uuid.length - 3)}
          </TooltipTrigger>
          <TooltipContent>
            <p>{uuid}</p>
          </TooltipContent>
        </Tooltip>
      </TooltipProvider>
    </div>
  );
};
