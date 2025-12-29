import csv


class CSVData:
    """docstring for CSVData."""

    def __init__(self, source_path, delimiter, io_type):
        self.source_path = source_path
        self.delimiter = delimiter

        if io_type == 'r':
            self.data = self.read_csv(self.source_path, self.delimiter)

            # Stores the headings apart from data
            self.headings = self.data[0]
            del self.data[0]
        elif io_type == 'w':
            self.data = []
        else:
            raise TypeError(
                'CSVData() must recieve an io_type value of \'r\' or \'w\''
                )

    # Reads "content" from "output_path" as a CSV file with "delimiter"
    def read_csv(self, source_path, delimiter):
        data = []
        with open(source_path, newline='') as csv_file:
            csv_reader = csv.reader(csv_file, delimiter=delimiter)
            for row in csv_reader:
                if(len(row) > 0):
                    if (row[0] != ''):
                        data.append(row)
        return data

    # Returns all data
    def get_data(self):
        return self.data

    # Returns the ids for all active students
    def get_active_stu_ids(self):
        active_stu_ids = []
        for stu in self.data:
            if 'Pending' not in stu[0]:
                active_stu_ids.append(int(stu[0]))

        return active_stu_ids

    def set_data(self, data):
        self.data = data
        return self

    # Writes "data" to "out_path" as a CSV file with "delimiter"
    def write_csv(self):
        with open(self.source_path, 'w', newline='') as csv_file:
            csv_writer = csv.writer(csv_file, delimiter=self.delimiter)
            for row in self.data:
                csv_writer.writerow(row)
        return self
