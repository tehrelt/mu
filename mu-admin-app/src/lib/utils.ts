import { clsx, type ClassValue } from "clsx";
import dayjs from "dayjs";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function datef(date: Date, format?: string): string {
  return dayjs(date).format(format ?? "DD/MM/YYYY HH:mm:ss");
}
