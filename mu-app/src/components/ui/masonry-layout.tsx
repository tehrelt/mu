// components/MasonryLayout.tsx
import { cn } from "@/shared/lib/utils";
import React from "react";

type MasonryLayoutProps = {
  children: React.ReactNode;
  className?: string;
};

const MasonryLayout: React.FC<MasonryLayoutProps> = ({
  children,
  className,
}) => {
  return (
    <div className={cn("columns-1 sm:columns-2 lg:columns-3 gap-4", className)}>
      {React.Children.map(children, (child, index) => (
        <div key={index} className="mb-4 break-inside-avoid max-w-full">
          {child}
        </div>
      ))}
    </div>
  );
};
export default MasonryLayout;
