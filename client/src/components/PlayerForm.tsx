import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Grid from "@mui/material/Grid";
import MenuItem from "@mui/material/MenuItem";
import useTheme from "@mui/material/styles/useTheme";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import React, { FormEvent, useState } from "react";

const raceOptions = [
  "Dwarf",
  "Gnome",
  "Human",
  "Night elf",
  "Orc",
  "Tauren",
  "Troll",
  "Undead",
];
const classOptions = [
  "Druid",
  "Hunter",
  "Mage",
  "Paladin",
  "Priest",
  "Rogue",
  "Shaman",
  "Warlock",
  "Warrior",
];
const professionOptions = [
  "Alchemy",
  "Blacksmithing",
  "Enchanting",
  "Engineering",
  "Herbalism",
  "Mining",
  "Tailoring",
];

interface PlayerFormProps {}

const PlayerForm: React.FunctionComponent<PlayerFormProps> = () => {
  const theme = useTheme();
  const [raceIndex, setRaceIndex] = useState<number | "">("");
  const [classIndex, setClassIndex] = useState<number | "">("");
  const [prof1Index, setProf1Index] = useState<number | "">("");
  const [prof2Index, setProf2Index] = useState<number | "">("");

  const selectStyle = {
    "& .MuiInputLabel-root": { color: "lightgrey" },
    "& .MuiInputLabel-root.Mui-focused": { color: "lightgrey" },
    "& .MuiOutlinedInput-root": {
      "& > fieldset": { borderColor: theme.palette.primary.main },
    },
    "& .MuiOutlinedInput-root.Mui-focused": {
      "& > fieldset": { borderColor: theme.palette.primary.main },
    },
    "& .MuiOutlinedInput-root:hover": {
      "& > fieldset": { borderColor: "lightgrey" },
    },
  };

  const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const f = new FormData(e.currentTarget);
    const resp = await fetch("/player", {
      method: "POST",
      body: new URLSearchParams(f as URLSearchParams),
    });
    console.log(resp);
  };

  const completed = () => raceIndex && !!classIndex;

  return (
    <Box component="form" noValidate autoComplete="off" onSubmit={handleSubmit}>
      <Grid container spacing={1}>
        <Grid item xs={6}>
          <TextField
            FormHelperTextProps={{ sx: { color: "white" } }}
            InputProps={{ sx: { color: "lightgrey" } }}
            sx={selectStyle}
            id="outlined-select-race"
            variant="outlined"
            select
            size="small"
            label="Race"
            value={raceIndex}
            onChange={(e) => setRaceIndex(e.target.value as unknown as number)}
            helperText="Select a Race"
          >
            {raceOptions.map((race, i) => (
              <MenuItem key={race} value={i}>
                <Typography>{race}</Typography>
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        <Grid item xs={6}>
          <TextField
            FormHelperTextProps={{ sx: { color: "white" } }}
            InputProps={{ sx: { color: "lightgrey" } }}
            sx={selectStyle}
            id="outlined-select-class"
            select
            size="small"
            label="Class"
            value={classIndex}
            onChange={(e) => setClassIndex(e.target.value as unknown as number)}
            helperText="Select a Class"
          >
            {classOptions.map((playerClass, i) => (
              <MenuItem key={playerClass} value={i}>
                <Typography>{playerClass}</Typography>
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        <Grid item xs={6}>
          <TextField
            FormHelperTextProps={{ sx: { color: "white" } }}
            InputProps={{ sx: { color: "lightgrey" } }}
            sx={selectStyle}
            id="outlined-select-prof1"
            variant="outlined"
            select
            size="small"
            label="Profession 1"
            value={prof1Index}
            onChange={(e) => setProf1Index(e.target.value as unknown as number)}
            helperText="Select a Profession"
          >
            {raceOptions.map((race, i) => (
              <MenuItem key={race} value={i}>
                <Typography>{race}</Typography>
              </MenuItem>
            ))}
          </TextField>
        </Grid>
        <Grid item xs={6}>
          <TextField
            FormHelperTextProps={{ sx: { color: "white" } }}
            InputProps={{ sx: { color: "lightgrey" } }}
            sx={selectStyle}
            id="outlined-select-prof2"
            select
            size="small"
            label="Profession 2"
            value={prof2Index}
            onChange={(e) => setProf2Index(e.target.value as unknown as number)}
            helperText="Select a Profession"
          >
            {classOptions.map((playerClass, i) => (
              <MenuItem key={playerClass} value={i}>
                <Typography>{playerClass}</Typography>
              </MenuItem>
            ))}
          </TextField>
        </Grid>
      </Grid>
      <Button
        disabled={!completed()}
        variant="contained"
        type="submit"
        color="primary"
      >
        Submit
      </Button>
    </Box>
  );
};

export default PlayerForm;
