import (
	"errors"
	"bufio"
	"unicode"
)

type parserState int

const (
	KEY_DEPRECATED		= 1 << iota
	KEY_LIST			= 1 << iota
	KEY_HIDDEN			= 1 << iota
	KEY_REMOVED			= 1 << iota
)

type keyDescriptor struct {
	key string
	flags int
	extra string
};

type sectionDescriptor struct {
	name string
	keys []keyDescriptor
	named bool
}

type section struct {
	name string
	descriptor sectionDescriptor	
}

type config struct {
	sections []section
}

var (
	keysSource []keyDescriptor = []keyDescriptor{
		{ "type",					0, "" },
		{ "sql_host",				0, "" },
		{ "sql_user",				0, "" },
		{ "sql_pass",				0, "" },
		{ "sql_db",					0, "" },
		{ "sql_port",				0, "" },
		{ "sql_sock",				0, "" },
		{ "mysql_connect_flags",	0, "" },
		{ "mysql_ssl_key",			0, "" },
		{ "mysql_ssl_cert",			0, "" },
		{ "mysql_ssl_ca",			0, "" },
		{ "mssql_winauth",			0, "" },
		{ "mssql_unicode",			KEY_REMOVED, "" },
		{ "sql_query_pre",			KEY_LIST, "" },
		{ "sql_query",				0, "" },
		{ "sql_query_range",		0, "" },
		{ "sql_range_step",			0, "" },
		{ "sql_query_killlist",		0, "" },
		{ "sql_attr_uint",			KEY_LIST, "" },
		{ "sql_attr_bool",			KEY_LIST, "" },
		{ "sql_attr_timestamp",		KEY_LIST, "" },
		{ "sql_attr_str2ordinal",	KEY_REMOVED | KEY_LIST, "" },
		{ "sql_attr_float",			KEY_LIST, "" },
		{ "sql_attr_bigint",		KEY_LIST, "" },
		{ "sql_attr_multi",			KEY_LIST, "" },
		{ "sql_query_post",			KEY_LIST, "" },
		{ "sql_query_post_index",	KEY_LIST, "" },
		{ "sql_ranged_throttle",	0, "" },
		{ "sql_query_info",			KEY_REMOVED, "" },
		{ "xmlpipe_command",		0, "" },
		{ "xmlpipe_field",			KEY_LIST, "" },
		{ "xmlpipe_attr_uint",		KEY_LIST, "" },
		{ "xmlpipe_attr_timestamp",	KEY_LIST, "" },
		{ "xmlpipe_attr_str2ordinal", KEY_REMOVED | KEY_LIST, "" },
		{ "xmlpipe_attr_bool",		KEY_LIST, "" },
		{ "xmlpipe_attr_float",		KEY_LIST, "" },
		{ "xmlpipe_attr_bigint",	KEY_LIST, "" },
		{ "xmlpipe_attr_multi",		KEY_LIST, "" },
		{ "xmlpipe_attr_multi_64",	KEY_LIST, "" },
		{ "xmlpipe_attr_string",	KEY_LIST, "" },
		{ "xmlpipe_attr_wordcount", KEY_REMOVED | KEY_LIST, "" },
		{ "xmlpipe_attr_json",		KEY_LIST, "" },
		{ "xmlpipe_field_string",	KEY_LIST, "" },
		{ "xmlpipe_field_wordcount", KEY_REMOVED | KEY_LIST, "" },
		{ "xmlpipe_fixup_utf8",		0, "" },
		{ "sql_str2ordinal_column", KEY_LIST | KEY_REMOVED, "" },
		{ "unpack_zlib",			KEY_LIST, "" },
		{ "unpack_mysqlcompress",	KEY_LIST, "" },
		{ "unpack_mysqlcompress_maxsize", 0, "" },
		{ "odbc_dsn",				0, "" },
		{ "sql_joined_field",		KEY_LIST, "" },
		{ "sql_attr_string",		KEY_LIST, "" },
		{ "sql_attr_str2wordcount", KEY_REMOVED | KEY_LIST, "" },
		{ "sql_field_string",		KEY_LIST, "" },
		{ "sql_field_str2wordcount", KEY_REMOVED | KEY_LIST, "" },
		{ "sql_file_field",			KEY_LIST, "" },
		{ "sql_column_buffers",		0, "" },
		{ "sql_attr_json",			KEY_LIST, "" },
		{ "hook_connect",			KEY_HIDDEN, "" },
		{ "hook_query_range",		KEY_HIDDEN, "" },
		{ "hook_post_index",		KEY_HIDDEN, "" },
		{ "tsvpipe_command",		0, "" },
		{ "tsvpipe_field",			KEY_LIST, "" },
		{ "tsvpipe_attr_uint",		KEY_LIST, "" },
		{ "tsvpipe_attr_timestamp",	KEY_LIST, "" },
		{ "tsvpipe_attr_bool",		KEY_LIST, "" },
		{ "tsvpipe_attr_float",		KEY_LIST, "" },
		{ "tsvpipe_attr_bigint",	KEY_LIST, "" },
		{ "tsvpipe_attr_multi",		KEY_LIST, "" },
		{ "tsvpipe_attr_multi_64",	KEY_LIST, "" },
		{ "tsvpipe_attr_string",	KEY_LIST, "" },
		{ "tsvpipe_attr_json",		KEY_LIST, "" },
		{ "tsvpipe_field_string",	KEY_LIST, "" },
		{ "csvpipe_command",		0, "" },
		{ "csvpipe_field",			KEY_LIST, "" },
		{ "csvpipe_attr_uint",		KEY_LIST, "" },
		{ "csvpipe_attr_timestamp",	KEY_LIST, "" },
		{ "csvpipe_attr_bool",		KEY_LIST, "" },
		{ "csvpipe_attr_float",		KEY_LIST, "" },
		{ "csvpipe_attr_bigint",	KEY_LIST, "" },
		{ "csvpipe_attr_multi",		KEY_LIST, "" },
		{ "csvpipe_attr_multi_64",	KEY_LIST, "" },
		{ "csvpipe_attr_string",	KEY_LIST, "" },
		{ "csvpipe_attr_json",		KEY_LIST, "" },
		{ "csvpipe_field_string",	KEY_LIST, "" },
		{ "csvpipe_delimiter",		0, "" },
	}

	keysIndex []keyDescriptor = []keyDescriptor{
		{ "source",					KEY_LIST, "" },
		{ "path",					0, "" },
		{ "docinfo",				0, "" },
		{ "mlock",					0, "" },
		{ "morphology",				0, "" },
		{ "stopwords",				0, "" },
		{ "exceptions",				0, "" },
		{ "wordforms",				KEY_LIST, "" },
		{ "embedded_limit",			0, "" },
		{ "min_word_len",			0, "" },
		{ "charset_type",			KEY_REMOVED, "" },
		{ "charset_table",			0, "" },
		{ "ignore_chars",			0, "" },
		{ "min_prefix_len",			0, "" },
		{ "min_infix_len",			0, "" },
		{ "max_substring_len",		0, "" },
		{ "prefix_fields",			0, "" },
		{ "infix_fields",			0, "" },
		{ "enable_star",			KEY_REMOVED, "" },
		{ "ngram_len",				0, "" },
		{ "ngram_chars",			0, "" },
		{ "phrase_boundary",		0, "" },
		{ "phrase_boundary_step",	0, "" },
		{ "ondisk_dict",			KEY_REMOVED, "" },
		{ "type",					0, "" },
		{ "local",					KEY_LIST, "" },
		{ "agent",					KEY_LIST, "" },
		{ "agent_blackhole",		KEY_LIST, "" },
		{ "agent_persistent",		KEY_LIST, "" },
		{ "agent_connect_timeout",	0, "" },
		{ "ha_strategy",			0, ""	},
		{ "agent_query_timeout",	0, "" },
		{ "html_strip",				0, "" },
		{ "html_index_attrs",		0, "" },
		{ "html_remove_elements",	0, "" },
		{ "preopen",				0, "" },
		{ "inplace_enable",			0, "" },
		{ "inplace_hit_gap",		0, "" },
		{ "inplace_docinfo_gap",	0, "" },
		{ "inplace_reloc_factor",	0, "" },
		{ "inplace_write_factor",	0, "" },
		{ "index_exact_words",		0, "" },
		{ "min_stemming_len",		0, "" },
		{ "overshort_step",			0, "" },
		{ "stopword_step",			0, "" },
		{ "blend_chars",			0, "" },
		{ "expand_keywords",		0, "" },
		{ "hitless_words",			0, "" },
		{ "hit_format",				KEY_HIDDEN | KEY_DEPRECATED, "default value" },
		{ "rt_field",				KEY_LIST, "" },
		{ "rt_attr_uint",			KEY_LIST, "" },
		{ "rt_attr_bigint",			KEY_LIST, "" },
		{ "rt_attr_float",			KEY_LIST, "" },
		{ "rt_attr_timestamp",		KEY_LIST, "" },
		{ "rt_attr_string",			KEY_LIST, "" },
		{ "rt_attr_multi",			KEY_LIST, "" },
		{ "rt_attr_multi_64",		KEY_LIST, "" },
		{ "rt_attr_json",			KEY_LIST, "" },
		{ "rt_attr_bool",			KEY_LIST, "" },
		{ "rt_mem_limit",			0, "" },
		{ "dict",					0, "" },
		{ "index_sp",				0, "" },
		{ "index_zones",			0, "" },
		{ "blend_mode",				0, "" },
		{ "regexp_filter",			KEY_LIST, "" },
		{ "bigram_freq_words",		0, "" },
		{ "bigram_index",			0, "" },
		{ "index_field_lengths",	0, "" },
		{ "divide_remote_ranges",	KEY_HIDDEN, "" },
		{ "stopwords_unstemmed",	0, "" },
		{ "global_idf",				0, "" },
		{ "rlp_context",			0, "" },
		{ "ondisk_attrs",			0, "" },
		{ "index_token_filter",		0, "" },
	}

	keysIndexer []keyDescriptor = []keyDescriptor{
		{ "mem_limit",				0, "" },
		{ "max_iops",				0, "" },
		{ "max_iosize",				0, "" },
		{ "max_xmlpipe2_field",		0, "" },
		{ "max_file_field_buffer",	0, "" },
		{ "write_buffer",			0, "" },
		{ "on_file_field_error",	0, "" },
		{ "on_json_attr_error",		KEY_DEPRECATED, "on_json_attr_error in common{..} section" },
		{ "json_autoconv_numbers",	KEY_DEPRECATED, "json_autoconv_numbers in common{..} section" },
		{ "json_autoconv_keynames",	KEY_DEPRECATED, "json_autoconv_keynames in common{..} section" },
		{ "lemmatizer_cache",		0, "" },
	}

	keysSearchd []keyDescriptor = []keyDescriptor{
		{ "address",				KEY_REMOVED, "" },
		{ "port",					KEY_REMOVED, "" },
		{ "listen",					KEY_LIST, "" },
		{ "log",					0, "" },
		{ "query_log",				0, "" },
		{ "read_timeout",			0, "" },
		{ "client_timeout",			0, "" },
		{ "max_children",			0, "" },
		{ "pid_file",				0, "" },
		{ "max_matches",			KEY_REMOVED, "" },
		{ "seamless_rotate",		0, "" },
		{ "preopen_indexes",		0, "" },
		{ "unlink_old",				0, "" },
		{ "ondisk_dict_default",	KEY_REMOVED, "" },
		{ "attr_flush_period",		0, "" },
		{ "max_packet_size",		0, "" },
		{ "mva_updates_pool",		0, "" },
		{ "max_filters",			0, "" },
		{ "max_filter_values",		0, "" },
		{ "listen_backlog",			0, "" },
		{ "read_buffer",			0, "" },
		{ "read_unhinted",			0, "" },
		{ "max_batch_queries",		0, "" },
		{ "subtree_docs_cache",		0, "" },
		{ "subtree_hits_cache",		0, "" },
		{ "workers",				0, "" },
		{ "prefork",				KEY_HIDDEN, "" },
		{ "dist_threads",			0, "" },
		{ "binlog_flush",			0, "" },
		{ "binlog_path",			0, "" },
		{ "binlog_max_log_size",	0, "" },
		{ "thread_stack",			0, "" },
		{ "expansion_limit",		0, "" },
		{ "rt_flush_period",		0, "" },
		{ "query_log_format",		0, "" },
		{ "mysql_version_string",	0, "" },
		{ "plugin_dir",				KEY_DEPRECATED, "plugin_dir in common{..} section" },
		{ "collation_server",		0, "" },
		{ "collation_libc_locale",	0, "" },
		{ "watchdog",				0, "" },
		{ "prefork_rotation_throttle", 0, "" },
		{ "snippets_file_prefix",	0, "" },
		{ "sphinxql_state",			0, "" },
		{ "rt_merge_iops",			0, "" },
		{ "rt_merge_maxiosize",		0, "" },
		{ "ha_ping_interval",		0, "" },
		{ "ha_period_karma",		0, "" },
		{ "predicted_time_costs",	0, "" },
		{ "persistent_connections_limit",	0, "" },
		{ "ondisk_attrs_default",	0, "" },
		{ "shutdown_timeout",		0, "" },
		{ "query_log_min_msec",		0, "" },
		{ "agent_connect_timeout",	0, "" },
		{ "agent_query_timeout",	0, "" },
		{ "agent_retry_delay",		0, "" },
		{ "agent_retry_count",		0, "" },
		{ "net_wait_tm",			0, "" },
		{ "net_throttle_action",	0, "" },
		{ "net_throttle_accept",	0, "" },
		{ "net_send_job",			0, "" },
		{ "net_workers",			0, "" },
		{ "queue_max_length",		0, "" },
	}

	keysCommon []keyDescriptor = []keyDescriptor{
		{ "lemmatizer_base",		0, "" },
		{ "on_json_attr_error",		0, "" },
		{ "json_autoconv_numbers",	0, "" },
		{ "json_autoconv_keynames",	0, "" },
		{ "rlp_root",				0, "" },
		{ "rlp_environment",		0, "" },
		{ "rlp_max_batch_size",		0, "" },
		{ "rlp_max_batch_docs",		0, "" },
		{ "plugin_dir",				0, "" },
	}

	configSections map[string]sectionDescriptor = map[string]sectionDescriptor{
		SECTION_SOURCE:		{ SECTION_SOURCE,	keysSource,		true },
		SECTION_INDEX:		{ SECTION_INDEX,  	keysIndex,		true },
		SECTION_INDEXER:	{ SECTION_INDEXER,	keysIndexer,	false },
		SECTION_SEARCHD:	{ SECTION_SEARCHD,	keysSearchd,	false },
		SECTION_COMMON: 	{ SECTION_COMMON,	keysCommon,		false },
	}
)

const (
	S_TOP parserState = iota
	S_SKIP2NL
	S_TOK
	S_TYPE
	S_SEC
	S_CHR
	S_VALUE
	S_SECNAME
	S_SECBASE
	S_KEY
)

const (
	maxTokenLength = 64
)

/*func addSection(cfg config, type string, name string) (section, error) {

	return nil, nil

	
	m_sSectionType = sType;
	m_sSectionName = sName;

	if ( !m_tConf.Exists ( m_sSectionType ) )
		m_tConf.Add ( CSphConfigType(), m_sSectionType ); // FIXME! be paranoid, verify that it returned true

	if ( m_tConf[m_sSectionType].Exists ( m_sSectionName ) )
	{
		snprintf ( m_sError, sizeof(m_sError), "section '%s' (type='%s') already exists", sName, sType );
		return false;
	}
	m_tConf[m_sSectionType].Add ( CSphConfigSection(), m_sSectionName ); // FIXME! be paranoid, verify that it returned true

	return true;	
}*/

func addSection(cfg config, _type, name string) (section, error) {
	for _, sec := range cfg.sections {
		if sec.type == _type {
			return nil, "The unnamed section already exists"
		}
		if sec.name == name {
			return nil, "The named section already exists"
		}		
	}

	return section{ descriptor: configSections[_type], name: name }
}

func isAlphaEx(c rune) bool {
	return (c >= '0' && c <= '9') || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == '-' || c == '_';
}

func isPlainSection(key string) bool {	
	if section := configSections[key], section != nil && !section.named {
		return true
	}

	return false
}

func isNamedSection(key string) bool {	
	if section := configSections[key], section != nil && section.named {
		return true
	}

	return false
}

func Parse(file *os.File) (string, error) {
	var state parserState = S_TOP
	var stateStack [8]parserState
	stackPos := 0
	charIdx := 0
	iToken := 0
	iCh := '\x00'
	push := func (st parserState) {
		// TODO: Boundary check
		stateStack[stackPos] = st
		state = st
		stackPos++
	}

	pop := func () {
		stackPos--
		state = stateStack[stackPos]
	}

	back := func() { charIdx-- }
	err := func (format string, a ... interface{}) (string, error) { return "", errors.New(fmt.Sprintf(format, a...)) }

	scanner := bufio.NewScanner(file)
    for scanner.Scan() {
    	line := scanner.Text()    	
		charIdx = 0

		for ; charIdx < len(line); charIdx++ {
			p := line[charIdx]
			if state == S_TOP {
				if unicode.IsSpace(p) {
					continue
				}

				if p == '#' {
					// TODO: If on Linux, check if we are on the first line and if it is indeed a hashbang. If it is, try to execute the config file using the specified shell.
					//
					push(S_SKIP2NL)	
					continue
				}

				if !isAlphaEx(p) {
					return err("invalid token")
				}

				iToken = 0
				push(S_TYPE)
				push(S_TOK)
				back()
				continue
    		}

    		if state == S_SKIP2NL {
    			pop()
    			charIdx = len(line)
    			continue
    		}

    		if state == S_TOK {
    			if iToken == 0 && !isAlphaEx(p) {
    				return err("internal error (non-alpha in S_TOK pos 0)")
    			}

    			if len(sToken) == maxTokenLength {
    				return err("token too long")
    			}

    			if !isAlphaEx(p) {
    				pop()
    				sToken = ""
    				iToken = 0
    				back()
    				continue
    			}

    			if iToken == 0 {
    				sToken = ""
    			}

    			sToken += p
    			iToken++
    			continue
    		}

    		if state == S_TYPE {		
				if unicode.IsSpace(p) {
					continue
				}

				if p == '#' {
					push(S_SKIP2NL)
					continue
				}

				if sToken == "" {
					return err("internal error (empty token in S_TYPE)")
				}

				if isPlainSection(sToken) {
					if !AddSection(sToken, sToken) {
						break
					}
					sToken = ""
					pop()
					push(S_SEC)
					push(S_CHR)
					iCh = '{'
					back()
					continue
				}
				if IsNamedSection(sToken) {
					m_sSectionType = sToken
					sToken = ""
					pop()
					push(S_SECNAME)
					back()
					continue
				}
				
				return err("invalid section type '%s'", sToken)
			}
    	}
    }
}