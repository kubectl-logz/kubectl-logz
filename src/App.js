import * as React from "react";
import Container from "@mui/material/Container";
import { useEffect, useState } from "react";
import {
  Alert,
  Box,
  Icon,
  MenuItem,
  Paper,
  Select,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableFooter,
  TableHead,
  TablePagination,
  TableRow,
} from "@mui/material";
import Moment from "react-moment";

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
  const [level, setLevel] = useState("debug");
  const [entries, setEntries] = useState();
  const [files, setFiles] = useState();
  const [error, setError] = useState();

  useEffect(() => {
    fetch("/api/v1/logs")
      .then((r) => r.json())
      .then((j) => setFiles(j.files))
      .catch(setError);
  }, []);

  useEffect(() => {
    setPage(0);
  }, [file]);

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
          limit
      )
        .then((r) => r.json())
        .then((j) => {
          setCount(j.metadata.count);
          setEntries(j.entries);
        })
        .catch(setError);
    }
  }, [file, level, page, limit]);

  if (error) {
    return <Alert severity={"error"}>{error.message}</Alert>;
  }

  return (
    <Box>
      <Select onChange={(e) => setFile(e.target.value)} value={file || ""}>
        <MenuItem>Select log...</MenuItem>
        {files?.map((item) => (
          <MenuItem key={item} value={item}>
            {item}
          </MenuItem>
        ))}
      </Select>
      <Select onChange={(e) => setLevel(e.target.value)} value={level}>
        <MenuItem key="error" value="error">
          error
        </MenuItem>
        <MenuItem key="warn" value="warn">
          warn
        </MenuItem>
        <MenuItem key="info" value="info">
          info
        </MenuItem>
        <MenuItem key="debug" value="debug">
          debug
        </MenuItem>
      </Select>
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
