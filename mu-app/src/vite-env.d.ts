/// <reference types="vite/client" />

declare global {
  namespace NodeJS {
    interface ProcessEnv {
      API_ADDRESS: string;
    }
  }
}
