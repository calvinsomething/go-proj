import React from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import { ThemeProvider, createTheme } from "@mui/material";

import "./App.css";

import Home from "./containers/Home";

function App() {
  const theme = createTheme({
    palette: {
      primary: {
        main: "#ffb813",
        dark: "#ffb813",
      },
      action: { disabledBackground: "#b49a5f", hover: "#e2d4b5" },
      secondary: {
        main: "#2dace0",
        dark: "#2dace0",
      },
      info: {
        main: "#c10000",
        dark: "#c10000",
      },
    },
  });

  return (
    <div className="App">
      <BrowserRouter>
        <ThemeProvider theme={theme}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="*" />
          </Routes>
        </ThemeProvider>
      </BrowserRouter>
    </div>
  );
}

export default App;
