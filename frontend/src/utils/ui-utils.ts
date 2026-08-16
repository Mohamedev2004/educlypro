import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

/**
 * UI utility functions.
 * 
 * Responsibility: Provide helper functions for styling and UI logic.
 * Layer: Utils
 */

/**
 * Combines tailwind classes with clsx and twMerge.
 */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}
