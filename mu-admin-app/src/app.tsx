import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import React from "react";
import { ThemeProvider } from "./app/providers/theme-provider";
import { BrowserRouter } from "react-router-dom";
import { RoutesConfig } from "./app/routes-config";
import { Toaster } from "./components/ui/sonner";

function App() {
  const [queryClient] = React.useState(
    new QueryClient({
      defaultOptions: {
        queries: {},
      },
    })
  );

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <BrowserRouter>
          <RoutesConfig />
        </BrowserRouter>
        <Toaster />
        {/* <ReactQueryDevtools initialIsOpen={false} /> */}
      </ThemeProvider>
    </QueryClientProvider>
  );
}

export default App;
