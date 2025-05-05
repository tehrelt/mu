import { useEffect, useState } from "react";

export const useTitle = (initial?: string) => {
  const [title, setTitle] = useState("MU");

  useEffect(() => {
    if (initial) {
      setTitle(initial);
    }
  }, [initial]);

  useEffect(() => {
    document.title = title;
  }, [title]);

  return { setTitle };
};
