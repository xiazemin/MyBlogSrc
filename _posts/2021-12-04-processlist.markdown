---
title: processlist
layout: post
category: mysql
author: 夏泽民
---
he status variables have the following meanings.

Aborted_clients

The number of connections that were aborted because the client died without closing the connection properly. See Section B.3.2.10, “Communication Errors and Aborted Connections”.

Aborted_connects

The number of failed attempts to connect to the MySQL server. See Section B.3.2.10, “Communication Errors and Aborted Connections”.

For additional connection-related information, check the Connection_errors_xxx status variables and the host_cache table.

Binlog_cache_disk_use

The number of transactions that used the temporary binary log cache but that exceeded the value of binlog_cache_size and used a temporary file to store statements from the transaction.

The number of nontransactional statements that caused the binary log transaction cache to be written to disk is tracked separately in the Binlog_stmt_cache_disk_use status variable.

Binlog_cache_use

The number of transactions that used the binary log cache.

Binlog_stmt_cache_disk_use

The number of nontransaction statements that used the binary log statement cache but that exceeded the value of binlog_stmt_cache_size and used a temporary file to store those statements.

Binlog_stmt_cache_use

The number of nontransactional statements that used the binary log statement cache.

Bytes_received

The number of bytes received from all clients.

Bytes_sent

The number of bytes sent to all clients.

Com_xxx

The Com_xxx statement counter variables indicate the number of times each xxx statement has been executed. There is one status variable for each type of statement. For example, Com_delete and Com_update count DELETE and UPDATE statements, respectively. Com_delete_multi and Com_update_multi are similar but apply to DELETE and UPDATE statements that use multiple-table syntax.

If a query result is returned from query cache, the server increments the Qcache_hits status variable, not Com_select. See Section 8.10.3.4, “Query Cache Status and Maintenance”.

All Com_stmt_xxx variables are increased even if a prepared statement argument is unknown or an error occurred during execution. In other words, their values correspond to the number of requests issued, not to the number of requests successfully completed.

The Com_stmt_xxx status variables are as follows:

Com_stmt_prepare

Com_stmt_execute

Com_stmt_fetch

Com_stmt_send_long_data

Com_stmt_reset

Com_stmt_close

Those variables stand for prepared statement commands. Their names refer to the COM_xxx command set used in the network layer. In other words, their values increase whenever prepared statement API calls such as mysql_stmt_prepare(), mysql_stmt_execute(), and so forth are executed. However, Com_stmt_prepare, Com_stmt_execute and Com_stmt_close also increase for PREPARE, EXECUTE, or DEALLOCATE PREPARE, respectively. Additionally, the values of the older statement counter variables Com_prepare_sql, Com_execute_sql, and Com_dealloc_sql increase for the PREPARE, EXECUTE, and DEALLOCATE PREPARE statements. Com_stmt_fetch stands for the total number of network round-trips issued when fetching from cursors.

Com_stmt_reprepare indicates the number of times statements were automatically reprepared by the server after metadata changes to tables or views referred to by the statement. A reprepare operation increments Com_stmt_reprepare, and also Com_stmt_prepare.

Compression

Whether the client connection uses compression in the client/server protocol.

Connection_errors_xxx

These variables provide information about errors that occur during the client connection process. They are global only and represent error counts aggregated across connections from all hosts. These variables track errors not accounted for by the host cache (see Section 5.1.11.2, “DNS Lookups and the Host Cache”), such as errors that are not associated with TCP connections, occur very early in the connection process (even before an IP address is known), or are not specific to any particular IP address (such as out-of-memory conditions).

Connection_errors_accept

The number of errors that occurred during calls to accept() on the listening port.

Connection_errors_internal

The number of connections refused due to internal errors in the server, such as failure to start a new thread or an out-of-memory condition.

Connection_errors_max_connections

The number of connections refused because the server max_connections limit was reached.

Connection_errors_peer_address

The number of errors that occurred while searching for connecting client IP addresses.

Connection_errors_select

The number of errors that occurred during calls to select() or poll() on the listening port. (Failure of this operation does not necessarily means a client connection was rejected.)

Connection_errors_tcpwrap

The number of connections refused by the libwrap library.

Connections

The number of connection attempts (successful or not) to the MySQL server.

Created_tmp_disk_tables

The number of internal on-disk temporary tables created by the server while executing statements.

You can compare the number of internal on-disk temporary tables created to the total number of internal temporary tables created by comparing Created_tmp_disk_tables and Created_tmp_tables values.

See also Section 8.4.4, “Internal Temporary Table Use in MySQL”.

Created_tmp_files

How many temporary files mysqld has created.

Created_tmp_tables

The number of internal temporary tables created by the server while executing statements.

You can compare the number of internal on-disk temporary tables created to the total number of internal temporary tables created by comparing Created_tmp_disk_tables and Created_tmp_tables values.

See also Section 8.4.4, “Internal Temporary Table Use in MySQL”.

Each invocation of the SHOW STATUS statement uses an internal temporary table and increments the global Created_tmp_tables value.

Delayed_errors

The number of rows written with INSERT DELAYED for which some error occurred (probably duplicate key).

This status variable is deprecated (because DELAYED inserts are deprecated); expect it to be removed in a future release.

Delayed_insert_threads

The number of INSERT DELAYED handler threads in use for nontransactional tables.

This status variable is deprecated (because DELAYED inserts are deprecated); expect it to be removed in a future release.

Delayed_writes

The number of INSERT DELAYED rows written to nontransactional tables.

This status variable is deprecated (because DELAYED inserts are deprecated); expect it to be removed in a future release.

Flush_commands

The number of times the server flushes tables, whether because a user executed a FLUSH TABLES statement or due to internal server operation. It is also incremented by receipt of a COM_REFRESH packet. This is in contrast to Com_flush, which indicates how many FLUSH statements have been executed, whether FLUSH TABLES, FLUSH LOGS, and so forth.

Handler_commit

The number of internal COMMIT statements.

Handler_delete

The number of times that rows have been deleted from tables.

Handler_external_lock

The server increments this variable for each call to its external_lock() function, which generally occurs at the beginning and end of access to a table instance. There might be differences among storage engines. This variable can be used, for example, to discover for a statement that accesses a partitioned table how many partitions were pruned before locking occurred: Check how much the counter increased for the statement, subtract 2 (2 calls for the table itself), then divide by 2 to get the number of partitions locked.

Handler_mrr_init

The number of times the server uses a storage engine's own Multi-Range Read implementation for table access.

Handler_prepare

A counter for the prepare phase of two-phase commit operations.

Handler_read_first

The number of times the first entry in an index was read. If this value is high, it suggests that the server is doing a lot of full index scans (for example, SELECT col1 FROM foo, assuming that col1 is indexed).

Handler_read_key

The number of requests to read a row based on a key. If this value is high, it is a good indication that your tables are properly indexed for your queries.

Handler_read_last

The number of requests to read the last key in an index. With ORDER BY, the server issues a first-key request followed by several next-key requests, whereas with ORDER BY DESC, the server issues a last-key request followed by several previous-key requests.

Handler_read_next

The number of requests to read the next row in key order. This value is incremented if you are querying an index column with a range constraint or if you are doing an index scan.

Handler_read_prev

The number of requests to read the previous row in key order. This read method is mainly used to optimize ORDER BY ... DESC.

Handler_read_rnd

The number of requests to read a row based on a fixed position. This value is high if you are doing a lot of queries that require sorting of the result. You probably have a lot of queries that require MySQL to scan entire tables or you have joins that do not use keys properly.

Handler_read_rnd_next

The number of requests to read the next row in the data file. This value is high if you are doing a lot of table scans. Generally this suggests that your tables are not properly indexed or that your queries are not written to take advantage of the indexes you have.

Handler_rollback

The number of requests for a storage engine to perform a rollback operation.

Handler_savepoint

The number of requests for a storage engine to place a savepoint.

Handler_savepoint_rollback

The number of requests for a storage engine to roll back to a savepoint.

Handler_update

The number of requests to update a row in a table.

Handler_write

The number of requests to insert a row in a table.

Innodb_available_undo_logs

The total number of available InnoDB rollback segments. Supplements the innodb_rollback_segments system variable, which defines the number of active rollback segments.

Innodb_buffer_pool_dump_status

The progress of an operation to record the pages held in the InnoDB buffer pool, triggered by the setting of innodb_buffer_pool_dump_at_shutdown or innodb_buffer_pool_dump_now.

For related information and examples, see Section 14.8.3.5, “Saving and Restoring the Buffer Pool State”.

Innodb_buffer_pool_load_status

The progress of an operation to warm up the InnoDB buffer pool by reading in a set of pages corresponding to an earlier point in time, triggered by the setting of innodb_buffer_pool_load_at_startup or innodb_buffer_pool_load_now. If the operation introduces too much overhead, you can cancel it by setting innodb_buffer_pool_load_abort.

For related information and examples, see Section 14.8.3.5, “Saving and Restoring the Buffer Pool State”.

Innodb_buffer_pool_bytes_data

The total number of bytes in the InnoDB buffer pool containing data. The number includes both dirty and clean pages. For more accurate memory usage calculations than with Innodb_buffer_pool_pages_data, when compressed tables cause the buffer pool to hold pages of different sizes.

Innodb_buffer_pool_pages_data

The number of pages in the InnoDB buffer pool containing data. The number includes both dirty and clean pages. When using compressed tables, the reported Innodb_buffer_pool_pages_data value may be larger than Innodb_buffer_pool_pages_total (Bug #59550).

Innodb_buffer_pool_bytes_dirty

The total current number of bytes held in dirty pages in the InnoDB buffer pool. For more accurate memory usage calculations than with Innodb_buffer_pool_pages_dirty, when compressed tables cause the buffer pool to hold pages of different sizes.

Innodb_buffer_pool_pages_dirty

The current number of dirty pages in the InnoDB buffer pool.

Innodb_buffer_pool_pages_flushed

The number of requests to flush pages from the InnoDB buffer pool.

Innodb_buffer_pool_pages_free

The number of free pages in the InnoDB buffer pool.

Innodb_buffer_pool_pages_latched

The number of latched pages in the InnoDB buffer pool. These are pages currently being read or written, or that cannot be flushed or removed for some other reason. Calculation of this variable is expensive, so it is available only when the UNIV_DEBUG system is defined at server build time.

Innodb_buffer_pool_pages_misc

The number of pages in the InnoDB buffer pool that are busy because they have been allocated for administrative overhead, such as row locks or the adaptive hash index. This value can also be calculated as Innodb_buffer_pool_pages_total − Innodb_buffer_pool_pages_free − Innodb_buffer_pool_pages_data. When using compressed tables, Innodb_buffer_pool_pages_misc may report an out-of-bounds value (Bug #59550).

Innodb_buffer_pool_pages_total

The total size of the InnoDB buffer pool, in pages. When using compressed tables, the reported Innodb_buffer_pool_pages_data value may be larger than Innodb_buffer_pool_pages_total (Bug #59550)

Innodb_buffer_pool_read_ahead

The number of pages read into the InnoDB buffer pool by the read-ahead background thread.

Innodb_buffer_pool_read_ahead_evicted

The number of pages read into the InnoDB buffer pool by the read-ahead background thread that were subsequently evicted without having been accessed by queries.

Innodb_buffer_pool_read_ahead_rnd

The number of “random” read-aheads initiated by InnoDB. This happens when a query scans a large portion of a table but in random order.

Innodb_buffer_pool_read_requests

The number of logical read requests.

Innodb_buffer_pool_reads

The number of logical reads that InnoDB could not satisfy from the buffer pool, and had to read directly from disk.

Innodb_buffer_pool_wait_free

Normally, writes to the InnoDB buffer pool happen in the background. When InnoDB needs to read or create a page and no clean pages are available, InnoDB flushes some dirty pages first and waits for that operation to finish. This counter counts instances of these waits. If innodb_buffer_pool_size has been set properly, this value should be small.

Innodb_buffer_pool_write_requests

The number of writes done to the InnoDB buffer pool.

Innodb_data_fsyncs

The number of fsync() operations so far. The frequency of fsync() calls is influenced by the setting of the innodb_flush_method configuration option.

Innodb_data_pending_fsyncs

The current number of pending fsync() operations. The frequency of fsync() calls is influenced by the setting of the innodb_flush_method configuration option.

Innodb_data_pending_reads

The current number of pending reads.

Innodb_data_pending_writes

The current number of pending writes.

Innodb_data_read

The amount of data read since the server was started (in bytes).

Innodb_data_reads

The total number of data reads (OS file reads).

Innodb_data_writes

The total number of data writes.

Innodb_data_written

The amount of data written so far, in bytes.

Innodb_dblwr_pages_written

The number of pages that have been written to the doublewrite buffer. See Section 14.12.1, “InnoDB Disk I/O”.

Innodb_dblwr_writes

The number of doublewrite operations that have been performed. See Section 14.12.1, “InnoDB Disk I/O”.

Innodb_have_atomic_builtins

Indicates whether the server was built with atomic instructions.

Innodb_log_waits

The number of times that the log buffer was too small and a wait was required for it to be flushed before continuing.

Innodb_log_write_requests

The number of write requests for the InnoDB redo log.

Innodb_log_writes

The number of physical writes to the InnoDB redo log file.

Innodb_num_open_files

The number of files InnoDB currently holds open.

Innodb_os_log_fsyncs

The number of fsync() writes done to the InnoDB redo log files.

Innodb_os_log_pending_fsyncs

The number of pending fsync() operations for the InnoDB redo log files.

Innodb_os_log_pending_writes

The number of pending writes to the InnoDB redo log files.

Innodb_os_log_written

The number of bytes written to the InnoDB redo log files.

Innodb_page_size

InnoDB page size (default 16KB). Many values are counted in pages; the page size enables them to be easily converted to bytes.

Innodb_pages_created

The number of pages created by operations on InnoDB tables.

Innodb_pages_read

The number of pages read from the InnoDB buffer pool by operations on InnoDB tables.

Innodb_pages_written

The number of pages written by operations on InnoDB tables.

Innodb_row_lock_current_waits

The number of row locks currently being waited for by operations on InnoDB tables.

Innodb_row_lock_time

The total time spent in acquiring row locks for InnoDB tables, in milliseconds.

Innodb_row_lock_time_avg

The average time to acquire a row lock for InnoDB tables, in milliseconds.

Innodb_row_lock_time_max

The maximum time to acquire a row lock for InnoDB tables, in milliseconds.

Innodb_row_lock_waits

The number of times operations on InnoDB tables had to wait for a row lock.

Innodb_rows_deleted

The number of rows deleted from InnoDB tables.

Innodb_rows_inserted

The number of rows inserted into InnoDB tables.

Innodb_rows_read

The number of rows read from InnoDB tables.

Innodb_rows_updated

The number of rows updated in InnoDB tables.

Innodb_truncated_status_writes

The number of times output from the SHOW ENGINE INNODB STATUS statement has been truncated.

Key_blocks_not_flushed

The number of key blocks in the MyISAM key cache that have changed but have not yet been flushed to disk.

Key_blocks_unused

The number of unused blocks in the MyISAM key cache. You can use this value to determine how much of the key cache is in use; see the discussion of key_buffer_size in Section 5.1.7, “Server System Variables”.

Key_blocks_used

The number of used blocks in the MyISAM key cache. This value is a high-water mark that indicates the maximum number of blocks that have ever been in use at one time.

Key_read_requests

The number of requests to read a key block from the MyISAM key cache.

Key_reads

The number of physical reads of a key block from disk into the MyISAM key cache. If Key_reads is large, then your key_buffer_size value is probably too small. The cache miss rate can be calculated as Key_reads/Key_read_requests.

Key_write_requests

The number of requests to write a key block to the MyISAM key cache.

Key_writes

The number of physical writes of a key block from the MyISAM key cache to disk.

Last_query_cost

The total cost of the last compiled query as computed by the query optimizer. This is useful for comparing the cost of different query plans for the same query. The default value of 0 means that no query has been compiled yet. The default value is 0. Last_query_cost has session scope.

Last_query_cost can be computed accurately only for simple, “flat” queries, but not for complex queries such as those containing subqueries or UNION. For the latter, the value is set to 0.

Last_query_partial_plans

The number of iterations the query optimizer made in execution plan construction for the previous query. Last_query_cost has session scope.

Max_used_connections

The maximum number of connections that have been in use simultaneously since the server started.

Not_flushed_delayed_rows

The number of rows waiting to be written to nontransactional tables in INSERT DELAYED queues.

This status variable is deprecated (because DELAYED inserts are deprecated); expect it to be removed in a future release.

Open_files

The number of files that are open. This count includes regular files opened by the server. It does not include other types of files such as sockets or pipes. Also, the count does not include files that storage engines open using their own internal functions rather than asking the server level to do so.

Open_streams

The number of streams that are open (used mainly for logging).

Open_table_definitions

The number of cached .frm files.

Open_tables

The number of tables that are open.

Opened_files

The number of files that have been opened with my_open() (a mysys library function). Parts of the server that open files without using this function do not increment the count.

Opened_table_definitions

The number of .frm files that have been cached.

Opened_tables

The number of tables that have been opened. If Opened_tables is big, your table_open_cache value is probably too small.

Performance_schema_xxx

Performance Schema status variables are listed in Section 22.16, “Performance Schema Status Variables”. These variables provide information about instrumentation that could not be loaded or created due to memory constraints.

Prepared_stmt_count

The current number of prepared statements. (The maximum number of statements is given by the max_prepared_stmt_count system variable.)

Qcache_free_blocks

The number of free memory blocks in the query cache.

Qcache_free_memory

The amount of free memory for the query cache.

Qcache_hits

The number of query cache hits.

The discussion at the beginning of this section indicates how to relate this statement-counting status variable to other such variables.

Qcache_inserts

The number of queries added to the query cache.

Qcache_lowmem_prunes

The number of queries that were deleted from the query cache because of low memory.

Qcache_not_cached

The number of noncached queries (not cacheable, or not cached due to the query_cache_type setting).

Qcache_queries_in_cache

The number of queries registered in the query cache.

Qcache_total_blocks

The total number of blocks in the query cache.

Queries

The number of statements executed by the server. This variable includes statements executed within stored programs, unlike the Questions variable. It does not count COM_PING or COM_STATISTICS commands.

The discussion at the beginning of this section indicates how to relate this statement-counting status variable to other such variables.

Questions

The number of statements executed by the server. This includes only statements sent to the server by clients and not statements executed within stored programs, unlike the Queries variable. This variable does not count COM_PING, COM_STATISTICS, COM_STMT_PREPARE, COM_STMT_CLOSE, or COM_STMT_RESET commands.

The discussion at the beginning of this section indicates how to relate this statement-counting status variable to other such variables.

Rpl_semi_sync_master_clients

The number of semisynchronous replicas.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_net_avg_wait_time

The average time in microseconds the source waited for a replica reply.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_net_wait_time

The total time in microseconds the source waited for replica replies.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_net_waits

The total number of times the source waited for replica replies.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_no_times

The number of times the source turned off semisynchronous replication.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_no_tx

The number of commits that were not acknowledged successfully by a replica.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_status

Whether semisynchronous replication currently is operational on the source. The value is ON if the plugin has been enabled and a commit acknowledgment has occurred. It is OFF if the plugin is not enabled or the source has fallen back to asynchronous replication due to commit acknowledgment timeout.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_timefunc_failures

The number of times the source failed when calling time functions such as gettimeofday().

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_tx_avg_wait_time

The average time in microseconds the source waited for each transaction.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_tx_wait_time

The total time in microseconds the source waited for transactions.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_tx_waits

The total number of times the source waited for transactions.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_wait_pos_backtraverse

The total number of times the source waited for an event with binary coordinates lower than events waited for previously. This can occur when the order in which transactions start waiting for a reply is different from the order in which their binary log events are written.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_wait_sessions

The number of sessions currently waiting for replica replies.

This variable is available only if the source-side semisynchronous replication plugin is installed.

Rpl_semi_sync_master_yes_tx

The number of commits that were acknowledged successfully by a replica.

This variable is available only if the master-side semisynchronous replication plugin is installed.

Rpl_semi_sync_slave_status

Whether semisynchronous replication currently is operational on the replica. This is ON if the plugin has been enabled and the replica I/O thread is running, OFF otherwise.

This variable is available only if the replica-side semisynchronous replication plugin is installed.

Rsa_public_key

This variable is available if MySQL was compiled using OpenSSL (see Section 6.3.4, “SSL Library-Dependent Capabilities”). Its value is the public key used by the sha256_password authentication plugin for RSA key pair-based password exchange. The value is nonempty only if the server successfully initializes the private and public keys in the files named by the sha256_password_private_key_path and sha256_password_public_key_path system variables. The value of Rsa_public_key comes from the latter file.

For information about sha256_password, see Section 6.4.1.4, “SHA-256 Pluggable Authentication”.

Select_full_join

The number of joins that perform table scans because they do not use indexes. If this value is not 0, you should carefully check the indexes of your tables.

Select_full_range_join

The number of joins that used a range search on a reference table.

Select_range

The number of joins that used ranges on the first table. This is normally not a critical issue even if the value is quite large.

Select_range_check

The number of joins without keys that check for key usage after each row. If this is not 0, you should carefully check the indexes of your tables.

Select_scan

The number of joins that did a full scan of the first table.

Slave_heartbeat_period

Shows the replication heartbeat interval (in seconds) on a replica.

Slave_last_heartbeat

Shows when the most recent heartbeat signal was received by a replica, as a TIMESTAMP value.

Slave_open_temp_tables

The number of temporary tables that the replica SQL thread currently has open. If the value is greater than zero, it is not safe to shut down the replica; see Section 17.4.1.29, “Replication and Temporary Tables”.

Slave_received_heartbeats

This counter increments with each replication heartbeat received by a replica since the last time that the replica was restarted or reset, or a CHANGE MASTER TO statement was issued.

Slave_retried_transactions

The total number of times since startup that the replica SQL thread has retried transactions.

Slave_rows_last_search_algorithm_used

The search algorithm that was most recently used by this replica to locate rows for row-based replication. The result shows whether the replica used indexes, a table scan, or hashing as the search algorithm for the last transaction executed on any channel.

The method used depends on the setting for the slave_rows_search_algorithms system variable, and the keys that are available on the relevant table.

This variable is available only for debug builds of MySQL.

Slave_running

This is ON if this server is a replica that is connected to a replication source, and both the I/O and SQL threads are running; otherwise, it is OFF.

Slow_launch_threads

The number of threads that have taken more than slow_launch_time seconds to create.

Slow_queries

The number of queries that have taken more than long_query_time seconds. This counter increments regardless of whether the slow query log is enabled. For information about that log, see Section 5.4.5, “The Slow Query Log”.

Sort_merge_passes

The number of merge passes that the sort algorithm has had to do. If this value is large, you should consider increasing the value of the sort_buffer_size system variable.

Sort_range

The number of sorts that were done using ranges.

Sort_rows

The number of sorted rows.

Sort_scan

The number of sorts that were done by scanning the table.

Ssl_accept_renegotiates

The number of negotiates needed to establish the connection.

Ssl_accepts

The number of accepted SSL connections.

Ssl_callback_cache_hits

The number of callback cache hits.

Ssl_cipher

The current encryption cipher (empty for unencrypted connections).

Ssl_cipher_list

The list of possible SSL ciphers (empty for non-SSL connections).

Ssl_client_connects

The number of SSL connection attempts to an SSL-enabled source.

Ssl_connect_renegotiates

The number of negotiates needed to establish the connection to an SSL-enabled source.

Ssl_ctx_verify_depth

The SSL context verification depth (how many certificates in the chain are tested).

Ssl_ctx_verify_mode

The SSL context verification mode.

Ssl_default_timeout

The default SSL timeout.

Ssl_finished_accepts

The number of successful SSL connections to the server.

Ssl_finished_connects

The number of successful replica connections to an SSL-enabled source.

Ssl_server_not_after

The last date for which the SSL certificate is valid. To check SSL certificate expiration information, use this statement:

mysql> SHOW STATUS LIKE 'Ssl_server_not%';
+-----------------------+--------------------------+
| Variable_name         | Value                    |
+-----------------------+--------------------------+
| Ssl_server_not_after  | Apr 28 14:16:39 2025 GMT |
| Ssl_server_not_before | May  1 14:16:39 2015 GMT |
+-----------------------+--------------------------+
In MySQL 5.6, the value is empty unless the connection uses SSL.

Ssl_server_not_before

The first date for which the SSL certificate is valid.

In MySQL 5.6, the value is empty unless the connection uses SSL.

Ssl_session_cache_hits

The number of SSL session cache hits.

Ssl_session_cache_misses

The number of SSL session cache misses.

Ssl_session_cache_mode

The SSL session cache mode.

Ssl_session_cache_overflows

The number of SSL session cache overflows.

Ssl_session_cache_size

The SSL session cache size.

Ssl_session_cache_timeouts

The number of SSL session cache timeouts.

Ssl_sessions_reused

How many SSL connections were reused from the cache.

Ssl_used_session_cache_entries

How many SSL session cache entries were used.

Ssl_verify_depth

The verification depth for replication SSL connections.

Ssl_verify_mode

The verification mode used by the server for a connection that uses SSL. The value is a bitmask; bits are defined in the openssl/ssl.h header file:

# define SSL_VERIFY_NONE                 0x00
# define SSL_VERIFY_PEER                 0x01
# define SSL_VERIFY_FAIL_IF_NO_PEER_CERT 0x02
# define SSL_VERIFY_CLIENT_ONCE          0x04
SSL_VERIFY_PEER indicates that the server asks for a client certificate. If the client supplies one, the server performs verification and proceeds only if verification is successful. SSL_VERIFY_CLIENT_ONCE indicates that a request for the client certificate is done only in the initial handshake.

Ssl_version

The SSL protocol version of the connection (for example, TLSv1). If the connection is not encrypted, the value is empty.

Table_locks_immediate

The number of times that a request for a table lock could be granted immediately.

Table_locks_waited

The number of times that a request for a table lock could not be granted immediately and a wait was needed. If this is high and you have performance problems, you should first optimize your queries, and then either split your table or tables or use replication.

Table_open_cache_hits

The number of hits for open tables cache lookups.

Table_open_cache_misses

The number of misses for open tables cache lookups.

Table_open_cache_overflows

The number of overflows for the open tables cache. This is the number of times, after a table is opened or closed, a cache instance has an unused entry and the size of the instance is larger than table_open_cache / table_open_cache_instances.

Tc_log_max_pages_used

For the memory-mapped implementation of the log that is used by mysqld when it acts as the transaction coordinator for recovery of internal XA transactions, this variable indicates the largest number of pages used for the log since the server started. If the product of Tc_log_max_pages_used and Tc_log_page_size is always significantly less than the log size, the size is larger than necessary and can be reduced. (The size is set by the --log-tc-size option. This variable is unused: It is unneeded for binary log-based recovery, and the memory-mapped recovery log method is not used unless the number of storage engines that are capable of two-phase commit and that support XA transactions is greater than one. (InnoDB is the only applicable engine.)

Tc_log_page_size

The page size used for the memory-mapped implementation of the XA recovery log. The default value is determined using getpagesize(). This variable is unused for the same reasons as described for Tc_log_max_pages_used.

Tc_log_page_waits

For the memory-mapped implementation of the recovery log, this variable increments each time the server was not able to commit a transaction and had to wait for a free page in the log. If this value is large, you might want to increase the log size (with the --log-tc-size option). For binary log-based recovery, this variable increments each time the binary log cannot be closed because there are two-phase commits in progress. (The close operation waits until all such transactions are finished.)

Threads_cached

The number of threads in the thread cache.

Threads_connected

The number of currently open connections.

Threads_created

The number of threads created to handle connections. If Threads_created is big, you may want to increase the thread_cache_size value. The cache miss rate can be calculated as Threads_created/Connections.

Threads_running

The number of threads that are not sleeping.

Uptime

The number of seconds that the server has been up.

Uptime_since_flush_status

The number of seconds since the most recent FLUSH STATUS statement.

https://dev.mysql.com/doc/refman/5.6/en/server-status-variables.html
<!-- more -->
mysql 支持三种连接方式

socket
named pipe
shared memory
named pipe 和 shared memory 只能在本地连接数据库，适用场景较少

thread_cache
参数 thread_cache_size 控制了 thread_cache 的大小， 设为0时关闭 thread_cache，不缓存空闲thread

mysql> show status like 'Threads%';
+-------------------+-------+
| Variable_name     | Value |
+-------------------+-------+
| Threads_cached    | 1     |
| Threads_connected | 1     |
| Threads_created   | 2     |
| Threads_running   | 1     |
+-------------------+-------+
4 rows in set (0.02 sec)

Threads_cached：缓存的 thread，新连接建立时，优先使用cache中的thread

Threads_connected：已连接的 thread

Threads_created：建立的 thread 数量

Threads_running：running状态的 thread 数量

Threads_created = Threads_cached + Threads_connected

Threads_running <= Threads_connected

MySQL 建立新连接非常消耗资源，频繁使用短连接，又没有其他组件实现连接池时，可以适当提高 thread_cache_size，降低新建连接的开销

.每个连接的限制
除了参数 max_user_connections 限制每个用户的最大连接数，还可以对每个用户制定更细致的限制

以下四个限制保存在mysql.user表中

MAX_QUERIES_PER_HOUR 每小时最大请求数（语句数量）
MAX_UPDATES_PER_HOUR 每小时最大更新数（更新语句的数量）
MAX_CONNECTIONS_PER_HOUR 每小时最大连接数
MAX_USER_CONNECTIONS 这个用户的最大连接数


http://mysql.taobao.org/monthly/2018/02/07/

客户购买的DB连接数是这个。max_connections，允许同时连接DB的客户端的最大线程数。如果客户端的连接数超过了max_connections,应用就会收到“too many connections”的错误。


已经创建的连接数
Threads_created是为处理连接而创建的线程数。再明确一点来说是连接到DB的，客户端的线程数。它包含Threads_running。 如果Threads_created很大，可能需要调整thread_cache_size。
线程cache命中率=Threads_created/Connections，cache命中率当然越大越好，如果命中率较低，可以考虑增加thread_cache_size。

https://developer.aliyun.com/article/683460


线程的状态信息：

已经创建的连接数
Threads_created是为处理连接而创建的线程数。再明确一点来说是连接到DB的，客户端的线程数。它包含Threads_running。 如果Threads_created很大，可能需要调整thread_cache_size。

线程cache命中率=Threads_created/Connections，cache命中率当然越大越好，如果命中率较低，可以考虑增加thread_cache_size。

已经连接的连接数
Thread_connected当前打开的连接数。

活跃连接数
Threads_running官方的说法是“没有sleep的线程数”。顾名思义是：在DB端正在执行的客户端线程总数。Server端保持这些连接同时客户端等待回复。有些线程可能消耗CPU或者IO，有些线程可能啥也没做单纯等表锁或行锁释放。当DB执行完这个线程，客户端收到回复，线程的状态就会从"running" 变成 "connected".

如果发现活跃链接数突然增高，通常是以下原因：

应用缓存失效
突发流量

https://cloud.tencent.com/developer/article/1816132

查看processlist这个表，表结构
ID：线程ID，这个信息对统计来说没有太大作用

USER：连接使用的账号，这个是一个统计维度，用于统计来自每个账号的连接数

HOST：连接客户端的IP/hostname+网络端口号，这也是一个统计维度，用于确定发起连接的客户端

DB：连接使用的default database，DB通常对应具体服务，可以用于判断服务的连接分布，这算一个统计维度

COMMAND：连接的动作，实际上是说连接处于哪个阶段，常见的有Sleep、Query、Connect、Statistics等，这也是一个统计维度，主要用于判断连接是否处于空闲状态

TIME：连接处于当前状态的时间，单位是s，这个在后面进行分析，暂不算在连接状态的统计维度中

STATE：连接的状态，表示当前MySQl连接正在做什么操作，这算一个统计维度，可能的值也比较多，详细可以查阅官方文档

INFO：连接正在执行的SQL，这个在下一节分析，暂不算在连接状态的统计维度中

https://dbaplus.cn/news-11-1396-1.html


mysql -uroot -h127.0.0.1 -e"use svc_t; show processlist;" |grep -v Sleep |wc -l
show status like 'Table%';
mysql -uroot -h127.0.0.1 -e"use svc_t;show status like 'Table%';"
mysql -uroot -h127.0.0.1 -e"use svc_t;show status like '%lock%';"
mysql -uroot -h127.0.0.1 -e"use svc_t;SELECT * FROM INFORMATION_SCHEMA.INNODB_LOCKS; "
mysql -uroot -h127.0.0.1 -e"use svc_t;SELECT * FROM INFORMATION_SCHEMA.INNODB_LOCK_WAITS;  "
mysql -uroot -h127.0.0.1 -e"use svc_t;show variables like '%timeout%'; "

后台线程
默认情况下，InnoDB 存储引擎有 13 个后台线程：

一个 master 线程

一个锁监控线程

一个错误监控线程

十个 IO 线程

插入缓存线程

日志线程

读线程（默认 4 个）

写线程（默认 4 个）

内存池
InnoDB 存储引擎的内存池包含：缓冲池、日志缓存池、额外内存池。这些内存的大小分别由配置文件中的参数决定。其中占比最大的是缓冲池，里面包含了数据缓存页、索引、插入缓存、自适应哈希索引、锁信息和数据字典。InnoDB 会在读取数据库数据的时候，将数据缓存到缓冲池中，而在修改数据的时候，会先把缓冲池中的数据修改掉，一旦修改过的数据页就会被标记为脏页，而脏页则会被 master 线程按照一定的频率刷新到磁盘中。日志缓存则是缓存了redo-log 信息，然后再刷新到 redo-log 文件中。额外内存池则是在对一些数据结构本身分配内存时会从额外内存池中申请内存，当该区域内存不足则会到缓冲池中申请。

Master Thread
InnoDB 存储引擎的主要工作都在一个单独的 Master Thread 中完成，其内部由四个循环体构成：主循环（ loop ）、后台循环（ background loop ）、刷新循环（ flush loop ）、暂停循环（ suspend loop ）
https://www.pianshen.com/article/476679761/

一、关于一个SQL的简单的工作过程
1、工作前提描述
　　1、启动MySQL，在内存中分配一个大空间innodb_buffer_pool(还有log_buffer)
　　2、多用户线程连接MySQL，从内存分配用户工作空间(其中排序空间)
　　3、磁盘上有数据库文件、ib_logfile、tmp目录、undo
2、SQL的简易流程
　　1、DQL操作
　　　　1、首先进行内存读
　　　　2、如果buffer pool中没有所需数据，就进行物理读
　　　　3、物理读数据读入buffer pool，再返回给用户工作空间
　　2、DML操作(例update)
　　　　1、内存读，然后进行物理读，读取所需修改的数据行
　　　　2、从磁盘调入undo页到buffer pool中
　　　　3、修改前的数据存入undo页里，产生redo
　　　　4、修改数据行(buffer pool中数据页成脏页)，产生redo
　　　　5、生成的redo先是存于用户工作空间，择机拷入log_buffer中
　　　　6、log线程不断的将log_buffer中的记录写入redo logfile中
　　　　7、修改完所有数据行，提交事务，刻意再触发一下log线程
　　　　8、待log_buffer中的相关信息都写完，响应事务提交成功
　　至此，日志写入磁盘，内存脏块还在buffer pool中(后台周期写入磁盘，释放buffer pool空间)。
　　
　　MySQL的工作机制是单进程多线程：IO线程=一个log线程+四个read线程+四个write线程
　　
　　1、读操作：innodb_read_io_threads
　　1、发起者：用户线程发起读请求
　　2、完成者：读线程执行请求队列中的读请求操作
　　3、如何调整读线程的数量
　　2、写操作：innodb_write_io_threads
　　1、发起者：page_cleaner线程发起
　　2、完成者：写线程执行请求队列中的写请求操作
　　3、如何调整写线程的数量
　　https://www.cnblogs.com/geaozhang/p/7214257.html
　　
　　
 Killing Threads (PROCESSLIST, KILL)
 Killing threads (KILL)
Once you've identified the problem thread, you can use the KILL command to kill it. There are basic two variations on the KILL command.

# Kill the entire connection.
KILL thread_id;
KILL CONNECTION thread_id;

# Terminate the currently executing statement, but leave the connection intact.
KILL QUERY thread_id;
https://oracle-base.com/articles/mysql/mysql-killing-threads

27.12.21.5 The processlist Table
ID

The connection identifier. This is the same value displayed in the Id column of the SHOW PROCESSLIST statement, displayed in the PROCESSLIST_ID column of the Performance Schema threads table, and returned by the CONNECTION_ID() function within the thread.

USER

The MySQL user who issued the statement. A value of system user refers to a nonclient thread spawned by the server to handle tasks internally, for example, a delayed-row handler thread or an I/O or SQL thread used on replica hosts. For system user, there is no host specified in the Host column. unauthenticated user refers to a thread that has become associated with a client connection but for which authentication of the client user has not yet occurred. event_scheduler refers to the thread that monitors scheduled events (see Section 25.4, “Using the Event Scheduler”).

Note
A USER value of system user is distinct from the SYSTEM_USER privilege. The former designates internal threads. The latter distinguishes the system user and regular user account categories (see Section 6.2.11, “Account Categories”).

HOST

The host name of the client issuing the statement (except for system user, for which there is no host). The host name for TCP/IP connections is reported in host_name:client_port format to make it easier to determine which client is doing what.

DB

The default database for the thread, or NULL if none has been selected.

COMMAND

The type of command the thread is executing on behalf of the client, or Sleep if the session is idle. For descriptions of thread commands, see Section 8.14, “Examining Server Thread (Process) Information”. The value of this column corresponds to the COM_xxx commands of the client/server protocol and Com_xxx status variables. See Section 5.1.10, “Server Status Variables”

TIME

The time in seconds that the thread has been in its current state. For a replica SQL thread, the value is the number of seconds between the timestamp of the last replicated event and the real time of the replica host. See Section 17.2.3, “Replication Threads”.

STATE

An action, event, or state that indicates what the thread is doing. For descriptions of STATE values, see Section 8.14, “Examining Server Thread (Process) Information”.

Most states correspond to very quick operations. If a thread stays in a given state for many seconds, there might be a problem that needs to be investigated.

INFO

The statement the thread is executing, or NULL if it is executing no statement. The statement might be the one sent to the server, or an innermost statement if the statement executes other statements. For example, if a CALL statement executes a stored procedure that is executing a SELECT statement, the INFO value shows the SELECT statement.

https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/performance-schema-processlist-table.html#function_connection-id
https://docs.oracle.com/cd/E17952_01/mysql-8.0-en/performance-schema-processlist-table.html

开启数据库的event执行调度

> 查看是否开启定时器

mysql> show variables like '%event_scheduler%';
+-----------------+-------+
| Variable_name   | Value |
+-----------------+-------+
| event_scheduler | OFF   |
+-----------------+-------+
https://www.cnblogs.com/geaozhang/p/6821692.html

How to share mysql connection between http goroutines?
The database/sql package manages the connection pooling automatically for you.

sql.Open(..) returns a handle which represents a connection pool, not a single connection. The database/sql package automatically opens a new connection if all connections in the pool are busy.

Applied to your code this means, that you just need to share the db-handle and use it in the HTTP handlers:

https://stackoverflow.com/questions/17376207/how-to-share-mysql-connection-between-http-goroutines

In Go 1.1 or newer, you can use db.SetMaxIdleConns(N) to limit the number of idle connections in the pool. This doesn’t limit the pool size, though.
In Go 1.2.1 or newer, you can use db.SetMaxOpenConns(N) to limit the number of total open connections to the database. Unfortunately, a deadlock bug (fix) prevents db.SetMaxOpenConns(N) from safely being used in 1.2.

http://go-database-sql.org/connection-pool.html
SetMaxOpenConns用于设置最大打开的连接数，默认值为0表示不限制。
SetMaxIdleConns用于设置闲置的连接数。

设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。


https://cloud.tencent.com/developer/article/1071721
http://hopehook.com/blog/golang_db_pool

MySQL的MaxIdleConns不合理，会变成短连接
是我们Go MySQL客户端最重要的配置。

maxIdleCount 最大空闲连接数，默认不配置，是2个最大空闲连接

maxOpen 最大连接数，默认不配置，是不限制最大连接数

maxLifetime 连接最大存活时间

maxIdleTime 空闲连接最大存活时间
当突发流量情况下，由于请求量级过大，超过了最大空闲连接数的负载，那么新的连接在放入连接池的时候，会被关闭，将连接变成短连接，导致服务性能进一步恶化。为了避免这种情况，下面列举了，可以优化的措施。

提前将maxIdleConns设大，避免出现短连接

做好mysql读写分离

提升mysql的吞吐量：精简返回字段，没必要的字段不要返回，能够够快复用连接

吞吐量的包尽量不要太大，避免分包

优化连接池，当客户端到MySQL的连接数大于最大空闲连接的时候，关闭能够做一下延迟（官方不支持，估计只能自己实现）

读请求的最好不要放MySQL里，尽量放redis里
https://blog.51cto.com/u_15127567/2714595
https://developpaper.com/golang-connection-pool-you-must-understand/
