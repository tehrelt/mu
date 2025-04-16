declare global {
  namespace NodeJS {
    interface ProcessEnv {
      API_ADDR: string;
      SESSION_SECRET: string;
    }
  }
}

export {};
