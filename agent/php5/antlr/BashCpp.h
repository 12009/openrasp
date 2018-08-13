
// Generated from BashCpp.g4 by ANTLR 4.7

#pragma once


#include "antlr4-runtime.h"




class  BashCpp : public antlr4::Lexer {
public:
  enum {
    COMMENT = 1, HERE_DOC = 2, SPACIAL_CAHR = 3, DOUBLE_CURLY = 4, VAR = 5, 
    ARGS = 6, MATH_EXEC = 7, EXEC = 8, DOLLAR = 9, FUNCTION_DEFINE = 10, 
    WS = 11, SINGLE_QUOTED_STR = 12, SINGLE_QUOTED_ERR = 13, DOUBLE_QUOTE = 14, 
    TICK = 15, UNQUOTED_STRING = 16, IN_DOUBLE_DOLLAR_TOKEN = 17, IN_DOUBLE_QUOTE_MATH_EXEC = 18, 
    IN_DOUBLE_EXEC = 19, TICK_START = 20, IN_DOUBLE_ERROR = 21, END_DOUBLE_QUOTE = 22, 
    DOUBLE_QUOTED_STRING = 23, IN_TICK_DOLLAR_TOKEN = 24, IN_TICK_COMMENT = 25, 
    IN_TICK_COMMENT2 = 26, IN_TICK_MATH_EXEC = 27, IN_TICK_EXEC = 28, IN_TICK_DOUBLE_QUOTE = 29, 
    END_TICK_TOKEN = 30, IN_TICK_SINGLE_QUOTED_ERR = 31, IN_TICK_HERE_DOC = 32, 
    IN_TICK_TOKEN = 33, IN_TICK_WS = 34, IN_EXEC_DOLLAR_TOKEN = 35, IN_EXEC_COMMENT = 36, 
    IN_EXEC_MATH_EXEC = 37, IN_EXEC_EXEC = 38, END_EXEC_TOKEN = 39, IN_EXEC_TICK_START = 40, 
    IN_EXEC_ERROR = 41, IN_EXEC_DOUBLE_QUOTE = 42, IN_EXEC_SINGLE_QUOTED_ERR = 43, 
    IN_EXEC_HERE_DOC = 44, IN_EXEC_TOKEN = 45, IN_EXEC_WS = 46, IN_TICK_IN_DOUBLE_DOLLAR_TOKEN = 47, 
    IN_TICK_IN_DOUBLE_QUOTE_MATH_EXEC = 48, IN_TICK_IN_DOUBLE_EXEC = 49, 
    IN_TICK_IN_DOUBLE_ERROR = 50, IN_TICK_END_DOUBLE_QUOTE = 51, IN_TICK_DOUBLE_QUOTED_STRING = 52, 
    IN_TICK_IN_EXEC_DOLLAR_TOKEN = 53, IN_TICK_IN_EXEC_COMMENT = 54, IN_TICK_IN_EXEC_MATH_EXEC = 55, 
    IN_TICK_IN_EXEC_EXEC = 56, IN_TICK_END_EXEC_TOKEN = 57, IN_TICK_IN_EXEC_ERROR = 58, 
    IN_TICK_IN_EXEC_DOUBLE_QUOTE = 59, IN_TICK_IN_EXEC_SINGLE_QUOTED_ERR = 60, 
    IN_TICK_IN_EXEC_HERE_DOC = 61, IN_TICK_IN_EXEC_TOKEN = 62, IN_TICK_IN_EXEC_WS = 63
  };

  enum {
    IN_DOUBLE_QUOTE = 1, IN_TICK = 2, IN_EXEC = 3, IN_TICK_IN_DOUBLE_QUOTE = 4, 
    IN_TICK_IN_EXEC = 5
  };

  BashCpp(antlr4::CharStream *input);
  ~BashCpp();


  char token_error = 0;

  virtual std::string getGrammarFileName() const override;
  virtual const std::vector<std::string>& getRuleNames() const override;

  virtual const std::vector<std::string>& getChannelNames() const override;
  virtual const std::vector<std::string>& getModeNames() const override;
  virtual const std::vector<std::string>& getTokenNames() const override; // deprecated, use vocabulary instead
  virtual antlr4::dfa::Vocabulary& getVocabulary() const override;

  virtual const std::vector<uint16_t> getSerializedATN() const override;
  virtual const antlr4::atn::ATN& getATN() const override;

  virtual void action(antlr4::RuleContext *context, size_t ruleIndex, size_t actionIndex) override;
private:
  static std::vector<antlr4::dfa::DFA> _decisionToDFA;
  static antlr4::atn::PredictionContextCache _sharedContextCache;
  static std::vector<std::string> _ruleNames;
  static std::vector<std::string> _tokenNames;
  static std::vector<std::string> _channelNames;
  static std::vector<std::string> _modeNames;

  static std::vector<std::string> _literalNames;
  static std::vector<std::string> _symbolicNames;
  static antlr4::dfa::Vocabulary _vocabulary;
  static antlr4::atn::ATN _atn;
  static std::vector<uint16_t> _serializedATN;


  // Individual action functions triggered by action() above.
  void COMMENTAction(antlr4::RuleContext *context, size_t actionIndex);
  void HERE_DOCAction(antlr4::RuleContext *context, size_t actionIndex);
  void EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void SINGLE_QUOTED_ERRAction(antlr4::RuleContext *context, size_t actionIndex);
  void DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void TICKAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_DOUBLE_QUOTE_MATH_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_DOUBLE_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void TICK_STARTAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_DOUBLE_ERRORAction(antlr4::RuleContext *context, size_t actionIndex);
  void END_DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_COMMENTAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_COMMENT2Action(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void END_TICK_TOKENAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_SINGLE_QUOTED_ERRAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_HERE_DOCAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_COMMENTAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void END_EXEC_TOKENAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_TICK_STARTAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_ERRORAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_SINGLE_QUOTED_ERRAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_EXEC_HERE_DOCAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_DOUBLE_QUOTE_MATH_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_DOUBLE_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_DOUBLE_ERRORAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_END_DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_COMMENTAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_EXECAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_END_EXEC_TOKENAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_ERRORAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_DOUBLE_QUOTEAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_SINGLE_QUOTED_ERRAction(antlr4::RuleContext *context, size_t actionIndex);
  void IN_TICK_IN_EXEC_HERE_DOCAction(antlr4::RuleContext *context, size_t actionIndex);

  // Individual semantic predicate functions triggered by sempred() above.

  struct Initializer {
    Initializer();
  };
  static Initializer _init;
};

