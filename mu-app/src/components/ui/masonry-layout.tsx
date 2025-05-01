// components/MasonryLayout.tsx
import React from "react";

type MasonryLayoutProps = {
  children: React.ReactNode;
};

const MasonryLayout: React.FC<MasonryLayoutProps> = ({ children }) => {
  return (
    <div className="columns-1 sm:columns-2 lg:columns-3 gap-4">
      {React.Children.map(children, (child, index) => (
        <div key={index} className="mb-4 break-inside-avoid">
          {child}
        </div>
      ))}
    </div>
  );
};
export default MasonryLayout;
