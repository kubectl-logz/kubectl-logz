import * as React from "react";
import Container from "@mui/material/Container";
import { useEffect, useState } from "react";
import dayjs from "dayjs";

import {
  Alert,
  Box,
  MenuItem,
  Select,
  Slider,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
  TextField,
  Toolbar,
} from "@mui/material";
import Moment from "react-moment";
import { DateTimePicker, LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";

export default function App() {
  return (
    <Container>
      <Logs />
    </Container>
  );
}

const Logs = () => {
  const [file, setFile] = useState();
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(10);
  const [count, setCount] = useState();
  const [level, setLevel] = useState("error");
  const [start, setStart] = useState(dayjs());
  const [entries, setEntries] = useState();
  const [files, setFiles] = useState();
  const [error, setError] = useState();

  useEffect(() => {
    fetch("/api/v1/logs")
      .then((r) => (r.ok ? r.json() : r.json().then((j) => Promise.reject(j))))
      .then((j) => setFiles(j.files))
      .catch(setError);
  }, []);

  useEffect(() => {
    setPage(0);
  }, [file]);

  useEffect(() => {
    if (error) {
      setTimeout(() => {
        setError(null);
      }, 10000);
    }
  }, [error]);

  useEffect(() => {
    if (file) {
      fetch(
        "/api/v1/logs/" +
          file +
          "?level=" +
          level +
          "&page=" +
          page +
          "&limit=" +
          limit +
          "&start=" +
          start.format()
      )
        .then((r) =>
          r.ok ? r.json() : r.json().then((j) => Promise.reject(j))
        )
        .then((j) => {
          setCount(j.metadata.count);
          setEntries(j.entries);
        })
        .catch(setError);
    }
  }, [file, level, page, limit, start]);

  const marks = [
    {
      value: 0,
      label: "debug",
    },
    {
      value: 1,
      label: "info",
    },
    {
      value: 2,
      label: "warn",
    },
    {
      value: 3,
      label: "error",
    },
  ];

  return (
    <Box>
      <Toolbar>
        <Box sx={{ width: 200, padding: 2 }}>
          <Slider
            value={marks.find((mark) => mark.label === level).value}
            max={3}
            step={1}
            valueLabelDisplay="auto"
            marks={marks}
            onChange={(event, newValue) => {
              setLevel(marks.find((mark) => mark.value === newValue).label);
            }}
          />
        </Box>
        <Box>
          <Select onChange={(e) => setFile(e.target.value)} value={file || ""}>
            <MenuItem>Select log...</MenuItem>
            {files?.map((item) => (
              <MenuItem key={item} value={item}>
                {item}
              </MenuItem>
            ))}
          </Select>
          <LocalizationProvider dateAdapter={AdapterDayjs}>
            <DateTimePicker
              onChange={(e) => {
                setStart(e);
              }}
              value={start}
              renderInput={(params) => <TextField {...params} />}
            />
          </LocalizationProvider>
        </Box>
      </Toolbar>
      {error && <Alert severity="error">{error.message}</Alert>}
      {entries && (
        <Box>
          <TableContainer>
            <Table size="small">
              <TableHead>
                <TableRow>
                  <TableCell>Time</TableCell>
                  <TableCell>Msg</TableCell>
                </TableRow>
              </TableHead>
              <TableBody>
                {entries?.map((entry, i) => (
                  <TableRow
                    key={i}
                    sx={{
                      backgroundColor: {
                        error: "#fee",
                        warning: "#ffe",
                      }[entry.level],
                    }}
                  >
                    <TableCell sx={{ color: "gray" }}>
                      {entry.time && (
                        <Moment ago durationFromNow trim>
                          {entry.time}
                        </Moment>
                      )}
                    </TableCell>
                    <TableCell>{entry.msg}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
              <TableFooter>
                <TableRow>
                  <TablePagination
                    count={count}
                    onPageChange={(event, newPage) => setPage(newPage)}
                    page={page}
                    rowsPerPage={limit}
                    onRowsPerPageChange={(event) =>
                      setLimit(parseInt(event.target.value, 10))
                    }
                    rowsPerPageOptions={[10, 50, 100, 500]}
                  />
                </TableRow>
              </TableFooter>
            </Table>
          </TableContainer>
        </Box>
      )}
    </Box>
  );
};
