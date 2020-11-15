from PyInstaller.utils.hooks import copy_metadata, collect_data_files
datas = copy_metadata('prettytable')
datas += collect_data_files('prettytable')
