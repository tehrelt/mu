import { themes, useTheme } from "@/app/providers/theme-provider";
import { cn } from "@/lib/utils";

const ThemeSwitcher = () => {
  const { theme, setTheme } = useTheme();

  return (
    <div className="px-2 py-1 flex justify-between bg-accent rounded-md">
      {Object.entries(themes.enum).map(([key, value]) => (
        <div
          key={key}
          className={cn(
            "px-4 py-1 cursor-pointer rounded-md",
            "hover:bg-background",
            theme === value && "bg-background"
          )}
          onClick={() => setTheme(value)}
        >
          {value}
        </div>
      ))}
    </div>
  );
};

export default ThemeSwitcher;
