import React, { useEffect, useState, useRef } from "react";
import Button from "@mui/material/Button";
import Popper from "@mui/material/Popper";

function Home() {
  const needsFetch = useRef(true);
  const [data, setData] = useState<any>(null);
  const [popper, setPopper] = useState<boolean>(false);
  const ref = useRef<HTMLButtonElement | null>(null);

  useEffect(() => {
    if (needsFetch.current) {
      needsFetch.current = false;
      fetch("/data")
        .then((resp) => resp.text())
        .then((data) => {
          setData(data);
        });
    }
  }, []);

  return (
    <>
      <Button
        variant="contained"
        color="secondary"
        onClick={() => setPopper(true)}
        ref={ref}
      >
        {data}
      </Button>
      <Popper anchorEl={ref.current} open={popper}>
        POPPER CONTENTS
      </Popper>
    </>
  );
}

export default Home;
