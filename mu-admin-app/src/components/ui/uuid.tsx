import React from "react";

type Props = {
  uuid: string;
  length?: number;
};

export const UUID = ({ uuid, length = 8 }: Props) => {
  return (
    <div>
      {uuid.substring(0, length)}...{uuid.substring(uuid.length - 3)}
    </div>
  );
};
