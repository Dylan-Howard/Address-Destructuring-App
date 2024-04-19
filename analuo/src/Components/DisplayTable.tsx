import { useState, ChangeEvent, MouseEvent } from 'react';
import TableContainer from '@mui/material/TableContainer';
import Paper from '@mui/material/Paper';
import TableHead from '@mui/material/TableHead';
import Table from '@mui/material/Table';
import TableRow from '@mui/material/TableRow';
import TableCell from '@mui/material/TableCell';
import TableBody from '@mui/material/TableBody';
import TablePagination from '@mui/material/TablePagination';
import { TablePaginationActionsProps } from '@mui/material/TablePagination/TablePaginationActions';
import IconButton from '@mui/material/IconButton';
import Box from '@mui/material/Box';

import LastPageIcon from '@mui/icons-material/LastPage';
import FirstPageIcon from '@mui/icons-material/FirstPage';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import Stack from '@mui/material/Stack';
import { SecondaryButton } from './Button';

function TablePaginationActions(props: TablePaginationActionsProps) {
  const { count, page, rowsPerPage, onPageChange } = props;

  const handleFirstPageButtonClick = (event: MouseEvent<HTMLButtonElement>) => {
    onPageChange(event, 0);
  }

  const handleBackButtonClick = (event: MouseEvent<HTMLButtonElement>) => {
    onPageChange(event, page - 1);
  };

  const handleNextButtonClick = (event: MouseEvent<HTMLButtonElement>) => {
    onPageChange(event, page + 1);
  };

  const handleLastPageButtonClick = (event: MouseEvent<HTMLButtonElement>) => {
    onPageChange(event, Math.max(0, Math.ceil(count / rowsPerPage) - 1));
  };

  return (
    <Box sx={{ flexShrink: 0, ml: 2.5 }}>
      <IconButton
        onClick={handleFirstPageButtonClick}
        disabled={page === 0}
        aria-label="first page"
      >
        <FirstPageIcon />
      </IconButton>
      <IconButton
        onClick={handleBackButtonClick}
        disabled={page === 0}
        aria-label="previous page"
      >
        <KeyboardArrowLeftIcon />
      </IconButton>
      <IconButton
        onClick={handleNextButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage) - 1}
        aria-label="next page"
      >
        <KeyboardArrowRightIcon />
      </IconButton>
      <IconButton
        onClick={handleLastPageButtonClick}
        disabled={page >= Math.ceil(count / rowsPerPage) - 1}
        aria-label="last page"
      >
        <LastPageIcon />
      </IconButton>
    </Box>
  );
}

function DisplayTable({ headers, rows, downloadUrl }: { headers: string[], rows: string[][], downloadUrl: string | undefined }) {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(+event.target.value);
    setPage(0);
  };

  if (rows.length < 2) {
    return (
      <Box>
        Hmm, it looks like nothing is here. 
      </Box>
    );
  }

  return (
    <Box>
      <TableContainer component={Paper} >
        <Table sx={{ display: 'block', maxWidth: 750, overflowX: 'scroll' }} size="small" aria-label="import data table">
          <TableHead>
            <TableRow>
              {
                headers.map((head) => <TableCell sx={{ fontWeight: 600 }}>{head}</TableCell>)
              }
            </TableRow>
          </TableHead>
          <TableBody>
            {rows
                .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                .map((row) => (
                  <TableRow
                    key={`${row[0]}-${row[1]}`}
                    sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
                  >
                    {
                      row.map((field) => <TableCell component="th" scope="row">{field}</TableCell>)
                    }
                  </TableRow>
                ))
            }
          </TableBody>
        </Table>
        <TablePagination
            rowsPerPageOptions={[10, 25, 100]}
            component="div"
            count={rows.length}
            rowsPerPage={rowsPerPage}
            page={page}
            onPageChange={handleChangePage}
            onRowsPerPageChange={handleChangeRowsPerPage}
            ActionsComponent={TablePaginationActions}
          />
      </TableContainer>
      {downloadUrl
        ? (
          <Stack direction="row" justifyContent="flex-end" sx={{ mt: 2 }}>
            <a href={`/data/export/${downloadUrl}`} download="Validation Results">
              <SecondaryButton variant="outlined">
                Download
              </SecondaryButton>
            </a>
          </Stack>
        )
        : ''
      }
    </Box>
  );
}

export default DisplayTable;
