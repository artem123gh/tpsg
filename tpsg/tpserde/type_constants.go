package tpserde

// Type layout - 32bit
//
// |      | not used | index type | structure type | scalar type / logical type |
// +------+----------+------------+--------+-----------+-----------------------------------------------+
// | size | 11 bits  |   4 bits   | 4 bits | 1 bit     |                  12 bits                      |
// |      |          |            |        |           +-----------+------------+-----------+----------+
// |      |          |            |        |           | 6 bits    |   4 bits   | 1 bit     | 1 bit    |
// |------+----------+------------+--------+-----------+-----------+------------+-----------+----------+
// |      |          |            |      TMASK         |                    SMASK                      |
// | mask |          |   IMASK    +--------+-----------+-----------+------------+-----------+----------+
// |      |          |            |        | T_CLONE   |           |   LMASK    | SC_CLONE  | SC_EVAL  |
// |      |          |            |        | type      |           |            | scalar    | scalar   |
// |      |          |            |        | clone bit |        STMASK          | clone bit | eval bit |

const (
    SC_ESHIFT = 1                    // scalar eval bit shift
    SC_CSHIFT = SC_ESHIFT + 1        // scalar clone bit shift
    SC_EVAL   = 1 << (SC_ESHIFT - 1) // bit to check if scalar type need special eval
    SC_CLONE  = 1 << (SC_CSHIFT - 1) // bit to check if scalar type need special cloning

    T_CSHIFT = SC_CSHIFT + 10        // type clone bit shift
    T_CLONE  = 1 << T_CSHIFT         // bit to check if type need special cloning

    PSHIFT = SC_CSHIFT + 4           // shift to get physical types from logical
    TSHIFT = T_CSHIFT + 1            // shift for structures/types
    LMASK  = 15 << SC_CSHIFT         // logical types mask
)

// SCALARS (12 bits), up to 15 logical types per each physical
const (
    NIL          uint32 = 0 << PSHIFT // must be 0
    SC_BOOL      uint32 = 1 << PSHIFT
    SC_BOOL_MAX  uint32 = SC_BOOL + LMASK
    SC_BYTE      uint32 = 2 << PSHIFT
    SC_BYTE_MAX  uint32 = SC_BYTE + LMASK
    SC_SHORT     uint32 = 3 << PSHIFT
    SC_SHORT_MAX uint32 = SC_SHORT + LMASK
    SC_INT       uint32 = 4 << PSHIFT
    SC_INT_MAX   uint32 = SC_INT + LMASK
    SC_MONTH     uint32 = SC_INT + (1 << SC_CSHIFT)
    SC_DATE      uint32 = SC_INT + (2 << SC_CSHIFT)
    SC_MINUTE    uint32 = SC_INT + (3 << SC_CSHIFT)
    SC_SECOND    uint32 = SC_INT + (4 << SC_CSHIFT)
    SC_TIME      uint32 = SC_INT + (5 << SC_CSHIFT)
    SC_LONG      uint32 = 5 << PSHIFT
    SC_LONG_MAX  uint32 = SC_LONG + LMASK
    SC_TIMESTAMP uint32 = SC_LONG + (1 << SC_CSHIFT)
    SC_DATETIME  uint32 = SC_LONG + (2 << SC_CSHIFT)
    SC_TIMESPAN  uint32 = SC_LONG + (3 << SC_CSHIFT)
    SC_INT128    uint32 = SC_CLONE | (6 << PSHIFT)
    SC_REAL      uint32 = 8 << PSHIFT
    SC_REAL_MAX  uint32 = SC_REAL + LMASK
    SC_FLOAT     uint32 = 9 << PSHIFT
    SC_FLOAT_MAX uint32 = SC_FLOAT + LMASK
    SC_QUAD      uint32 = SC_CLONE | (10 << PSHIFT)
    SC_GUID      uint32 = SC_CLONE | (12 << PSHIFT)
    SC_GUID_MAX  uint32 = SC_GUID + LMASK
    SC_ENUM      uint32 = 15 << PSHIFT
    SC_ENUM_MAX  uint32 = SC_ENUM + LMASK
    SC_SYMBOL    uint32 = (20 << PSHIFT) | SC_EVAL
    SC_ALIAS     uint32 = SC_SYMBOL + (1 << SC_CSHIFT)
    SC_SHADOW    uint32 = SC_SYMBOL + (2 << SC_CSHIFT)
    SC_SYMBOL_MAX uint32 = SC_SYMBOL + LMASK
    SC_CHAR      uint32 = 21 << PSHIFT
    SC_UPVAL     uint32 = (29 << PSHIFT) | SC_EVAL
    SC_REFTYPE   uint32 = (30 << PSHIFT) | SC_EVAL
    SC_REF_APPLY uint32 = SC_REFTYPE + (1 << SC_CSHIFT)
    SC_STACKTYPE uint32 = (31 << PSHIFT) | SC_EVAL
    SC_AST       uint32 = SC_CLONE | SC_EVAL | (32 << PSHIFT)
    SC_EXPR      uint32 = SC_AST + (1 << SC_CSHIFT)
    SMASK        uint32 = T_CLONE - 1
    STMASK       uint32 = SMASK ^ (SC_CLONE | SC_EVAL)
)

// STRUCTURES/TYPES (6 bits)
const (
    TMASK          uint32 = 31 << T_CSHIFT
    TSCALAR_MAX_SC uint32 = SC_AST - (1 << PSHIFT)
    TSCALAR        uint32 = 0 << TSHIFT
    TRAWPTR        uint32 = 1 << TSHIFT
    TCOPY          uint32 = 4 << TSHIFT
    TVEC           uint32 = T_CLONE | (5 << TSHIFT)
    TVEC_MAX_SC    uint32 = TVEC + SC_AST - (1 << PSHIFT)
    TVEC_MAX       uint32 = TVEC + SMASK
    TDEQUE         uint32 = T_CLONE | (7 << TSHIFT)
    TDEQUE_MAX     uint32 = TDEQUE + SMASK
    TOTHER         uint32 = T_CLONE | (15 << TSHIFT)
)

// INDICES (4 bits - max 15 index types)
const (
    ISHIFT         = TSHIFT + 4
    IMASK   uint32 = 15 << ISHIFT
    IDX_NONE uint32 = 0 << ISHIFT
    IDX_ASC  uint32 = 1 << ISHIFT
    IDX_DESC uint32 = 2 << ISHIFT
    IDX_SKIPLIST uint32 = 3 << ISHIFT
    IDX_MAX  uint32 = IDX_SKIPLIST
)

// VECTORS
const (
    VEC_BOOL      uint32 = SC_BOOL | TVEC
    VEC_BOOL_MAX  uint32 = VEC_BOOL + LMASK
    VEC_BYTE      uint32 = SC_BYTE | TVEC
    VEC_BYTE_MAX  uint32 = VEC_BYTE + LMASK
    VEC_SHORT     uint32 = SC_SHORT | TVEC
    VEC_SHORT_MAX uint32 = VEC_SHORT + LMASK
    VEC_INT       uint32 = SC_INT | TVEC
    VEC_INT_MAX   uint32 = VEC_INT + LMASK
    VEC_MONTH     uint32 = SC_MONTH | TVEC
    VEC_DATE      uint32 = SC_DATE | TVEC
    VEC_MINUTE    uint32 = SC_MINUTE | TVEC
    VEC_SECOND    uint32 = SC_SECOND | TVEC
    VEC_TIME      uint32 = SC_TIME | TVEC
    VEC_LONG      uint32 = SC_LONG | TVEC
    VEC_LONG_MAX  uint32 = VEC_LONG + LMASK
    VEC_TIMESTAMP uint32 = SC_TIMESTAMP | TVEC
    VEC_DATETIME  uint32 = SC_DATETIME | TVEC
    VEC_TIMESPAN  uint32 = SC_TIMESPAN | TVEC
    VEC_INT128    uint32 = SC_INT128 | TVEC
    VEC_GUID      uint32 = SC_GUID | TVEC
    VEC_GUID_MAX  uint32 = VEC_GUID + LMASK
    VEC_REAL      uint32 = SC_REAL | TVEC
    VEC_REAL_MAX  uint32 = VEC_REAL + LMASK
    VEC_FLOAT     uint32 = SC_FLOAT | TVEC
    VEC_FLOAT_MAX uint32 = VEC_FLOAT + LMASK
    VEC_QUAD      uint32 = SC_QUAD | TVEC
    VEC_REFTYPE   uint32 = SC_REFTYPE | TVEC
    VEC_SYMBOL    uint32 = SC_SYMBOL | TVEC
    VEC_ALIAS     uint32 = SC_ALIAS | TVEC
    VEC_SHADOW    uint32 = SC_SHADOW | TVEC
    VEC_SYMBOL_MAX uint32 = VEC_SYMBOL + LMASK
    VEC_CHAR      uint32 = SC_CHAR | TVEC
    VEC_ENUM      uint32 = SC_ENUM | TVEC
    VEC_ENUM_MAX  uint32 = VEC_ENUM + LMASK
    VEC_UPVAL     uint32 = SC_UPVAL | TVEC
    VEC_STACKTYPE uint32 = SC_STACKTYPE | TVEC
    LIST          uint32 = SC_AST | TVEC
    LIST_EXPR     uint32 = SC_EXPR | TVEC
)

// DEQUES
const (
    DEQ_BOOL      uint32 = SC_BOOL | TDEQUE
    DEQ_BOOL_MAX  uint32 = DEQ_BOOL + LMASK
    DEQ_BYTE      uint32 = SC_BYTE | TDEQUE
    DEQ_BYTE_MAX  uint32 = DEQ_BYTE + LMASK
    DEQ_SHORT     uint32 = SC_SHORT | TDEQUE
    DEQ_SHORT_MAX uint32 = DEQ_SHORT + LMASK
    DEQ_INT       uint32 = SC_INT | TDEQUE
    DEQ_INT_MAX   uint32 = DEQ_INT + LMASK
    DEQ_MONTH     uint32 = SC_MONTH | TDEQUE
    DEQ_DATE      uint32 = SC_DATE | TDEQUE
    DEQ_MINUTE    uint32 = SC_MINUTE | TDEQUE
    DEQ_SECOND    uint32 = SC_SECOND | TDEQUE
    DEQ_TIME      uint32 = SC_TIME | TDEQUE
    DEQ_LONG      uint32 = SC_LONG | TDEQUE
    DEQ_LONG_MAX  uint32 = DEQ_LONG + LMASK
    DEQ_TIMESTAMP uint32 = SC_TIMESTAMP | TDEQUE
    DEQ_DATETIME  uint32 = SC_DATETIME | TDEQUE
    DEQ_TIMESPAN  uint32 = SC_TIMESPAN | TDEQUE
    DEQ_INT128    uint32 = SC_INT128 | TDEQUE
    DEQ_GUID      uint32 = SC_GUID | TDEQUE
    DEQ_GUID_MAX  uint32 = DEQ_GUID + LMASK
    DEQ_REAL      uint32 = SC_REAL | TDEQUE
    DEQ_REAL_MAX  uint32 = DEQ_REAL + LMASK
    DEQ_FLOAT     uint32 = SC_FLOAT | TDEQUE
    DEQ_FLOAT_MAX uint32 = DEQ_FLOAT + LMASK
    DEQ_QUAD      uint32 = SC_QUAD | TDEQUE
    DEQ_REFTYPE   uint32 = SC_REFTYPE | TDEQUE
    DEQ_SYMBOL    uint32 = SC_SYMBOL | TDEQUE
    DEQ_SYMBOL_MAX uint32 = DEQ_SYMBOL + LMASK
    DEQ_CHAR      uint32 = SC_CHAR | TDEQUE
    DEQ_ENUM      uint32 = SC_ENUM | TDEQUE
    DEQ_ENUM_MAX  uint32 = DEQ_ENUM + LMASK
    DEQ_UPVAL     uint32 = SC_UPVAL | TDEQUE
    DEQ_AST       uint32 = SC_AST | TDEQUE
)

// EXECUTABLE/COPY
const (
    MONAD       uint32 = (0 << SC_CSHIFT) | TCOPY
    DYAD        uint32 = (1 << SC_CSHIFT) | TCOPY
    TRIAD       uint32 = (2 << SC_CSHIFT) | TCOPY
    TETRAD      uint32 = (3 << SC_CSHIFT) | TCOPY
    POLYAD      uint32 = (4 << SC_CSHIFT) | TCOPY
    COMMUTE     uint32 = (5 << SC_CSHIFT) | TCOPY
    RETURN      uint32 = (6 << SC_CSHIFT) | TCOPY | SC_EVAL
    ANY         uint32 = (7 << SC_CSHIFT) | TCOPY
    LAMBDA_REC  uint32 = (8 << SC_CSHIFT) | TCOPY | SC_EVAL
    TABLE_REF   uint32 = (9 << SC_CSHIFT) | TCOPY
    FIELD_REF   uint32 = (10 << SC_CSHIFT) | TCOPY
    AST_TYPE    uint32 = (11 << SC_CSHIFT) | TCOPY
    PROJECTION  uint32 = (12 << SC_CSHIFT) | TCOPY
)

// OTHER
const (
    LAMBDA       uint32 = (0 << SC_CSHIFT) | TOTHER
    REAGENT      uint32 = (1 << SC_CSHIFT) | TOTHER
    PATTERN      uint32 = (2 << SC_CSHIFT) | TOTHER | SC_EVAL
    TABLE        uint32 = (3 << SC_CSHIFT) | TOTHER
    DICT         uint32 = (4 << SC_CSHIFT) | TOTHER
    SELECT       uint32 = (6 << SC_CSHIFT) | TOTHER
    SELECT_C     uint32 = (7 << SC_CSHIFT) | TOTHER
    JOIN         uint32 = (8 << SC_CSHIFT) | TOTHER
    LJOIN        uint32 = (9 << SC_CSHIFT) | TOTHER
    CLOSURE      uint32 = (10 << SC_CSHIFT) | TOTHER | SC_EVAL
    TABLE_IDX    uint32 = (11 << SC_CSHIFT) | TOTHER
    DICT_TABLE   uint32 = (12 << SC_CSHIFT) | TOTHER
    PARSER       uint32 = (13 << SC_CSHIFT) | TOTHER
    USERDATA     uint32 = (14 << SC_CSHIFT) | TOTHER
    LAMBDA_WEAK  uint32 = (15 << SC_CSHIFT) | TOTHER | SC_EVAL
    LIST_WEAK    uint32 = (16 << SC_CSHIFT) | TOTHER | SC_EVAL
    TRACE        uint32 = (STMASK - (3 << SC_CSHIFT)) | TOTHER
    ERROR        uint32 = (18 << SC_CSHIFT) | TOTHER
    BREAKPOINT   uint32 = (STMASK - (2 << SC_CSHIFT)) | TOTHER
    INVALID      uint32 = (STMASK - (1 << SC_CSHIFT)) | TOTHER
)
