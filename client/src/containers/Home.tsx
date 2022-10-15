import React, { useEffect, useState, useRef } from "react";
import Typography from "@mui/material/Typography";

import PlayerForm from "../components/PlayerForm";
import Grid from "@mui/material/Grid";

function Home() {
  const loaded = useRef(false);
  const [playersData, setPlayersData] = useState<any>(null);

  useEffect(() => {
    if (!loaded.current) {
      loaded.current = true;
      fetch("/players")
        .then((resp) => resp.text())
        .then((data) => {
          setPlayersData(data);
        });
    }
  }, []);

  return (
    <Grid
      container
      sx={{ width: "90%", border: "solid 2px orange", p: 2, borderRadius: 4 }}
    >
      <Grid item xs={12} sm={8}>
        <Typography>{playersData}</Typography>
      </Grid>
      <Grid item xs={12} sm={4}>
        <PlayerForm />
      </Grid>
    </Grid>
  );
}

export default Home;
