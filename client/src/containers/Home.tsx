import React, { useEffect, useState } from "react";
import Button from "@mui/material/Button";

function Home() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    fetch("/data")
      .then((resp) => resp.text())
      .then((data) => {
        setData(data);
      });
  }, []);

  return <Button variant="contained">{data}</Button>;
}

export default Home;
