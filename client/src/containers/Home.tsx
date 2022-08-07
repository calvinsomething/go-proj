import React, { useEffect, useState } from "react";

function Home() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    fetch("/data")
      .then((resp) => resp.text())
      .then((data) => {
        setData(data);
      });
  }, []);

  return <p>{data}</p>;
}

export default Home;
