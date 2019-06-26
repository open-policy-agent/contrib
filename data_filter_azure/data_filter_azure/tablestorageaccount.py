from azure.storage import SharedAccessSignature
from azure.storage.table import TableService, Entity

class TableStorageAccount(object):
    """
    Provides a factory for creating table
    with a common account name and connection_string. 
    """

    def __init__(self, account_name=None, connection_string=None, sas_token=None, endpoint_suffix = 'cosmosdb.windows.net', is_emulated=None):
        '''
        :param str account_name:
            Storage account account name.
        :param str connection_string:
            Storage account connection string.
        :param str sas_token:
            Storage account sas token.
        :param str enpoint_suffix:
            Storage account endpoint_suffix.
        :param bool is_emulated:
            Whether to use the emulator. Defaults to False. If specified, will 
            override all other parameters.
        '''
        self.account_name = account_name
        self.connection_string = connection_string
        self.sas_token = sas_token
        self.endpoint_suffix = endpoint_suffix		
        self.is_emulated = is_emulated

    def create_table_service(self):
        '''
        Creates a TableService object with the settings specified in the 
        TableStorageAccount.

        :return: A service object.
        :rtype: :class:`~azure.storage.table.tableservice.TableService`
        '''
        return TableService(account_name = self.account_name,
                            sas_token=self.sas_token,
                            endpoint_suffix=self.endpoint_suffix, 
                            connection_string= self.connection_string,
                            is_emulated=self.is_emulated)

    def is_azure_cosmosdb_table(self):
        return self.connection_string != None and "table.cosmosdb" in self.connection_string