import os
from .CSVData import CSVData
from .AddressDatabase import AddressDatabase

unit_col_in = 5
notes_col_in = 16


# Writes the data to the specified, comma-delimited CSV file
def csv_quick_write(dir, filename, data):
    csv = CSVData(os.path.join(dir, filename + '.csv'), ',', 'w')
    csv.set_data(data).write_csv()
    # csv.write_csv()
    return filename


# Searches specificed directry and returns either a matching file or None
def get_file(dir, filenameStart, extension):
    for file in os.listdir(dir):
        if file.startswith(filenameStart) and file.endswith(extension):
            return file

    return None


# Transforms addresses and writes
def transform_addresses(data_dir, data_file, out_dir):
    data_file = get_file(data_dir, data_file, '.csv')

    print(data_file)

    if data_file is None:  # Ends script if no import data exists
        print('No import data could be found.')
        quit()

    addr_data = CSVData(os.path.join(data_dir, data_file), ',', 'r').get_data()
    address_db = AddressDatabase(addr_data, unit_col_in, notes_col_in)

    # Writes all formatted addresses to one file
    csv_quick_write(out_dir, 'addr_output', address_db.get_formatted())

    # Writes all non-WCPS addresses to one file
    csv_quick_write(out_dir, 'not_wcps', address_db.get_not_wcps())

    # Writes all addresses needing manual transformations to one file
    csv_quick_write(out_dir, 'needs_transforms', address_db.get_manual_transforms())
