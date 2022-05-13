#import "stdlib/core.b"

string root_source;

void init(){
	get_root_source(root_source);
	SET_SOURCE(root_source);
	SET_DEST("benv/print_format_program.b");
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

stack ops(string command, string op){
	stack quotes;
	stack these_quotes; 
	stack op_nums;
	string buf;
	stack res;
	int num1;
	int num2; 
	int op_num;
	bool was_quote;
	bool to_add;
	
	was_quote = False;
	to_add = True;

	op_nums = indexes(command, op);
	op_nums.pop(buf);
	#op_nums_s:
	[goto(#op_nums_e), ("end" == buf), print("")];
	op_num = int(buf);
	quotes = reg_find("\"(\\.|[^\"])*\"", command);
	
	#quotes_s:
	quotes.pop(these_quotes);
	these_quotes.pop(buf);
	[goto(#quotes_e), ("end" == buf), print("")];
	#these_quotes_s:
	[goto(#these_quotes_e), ("end" == buf), print("")];
	num1 = int(buf);
	these_quotes.pop(buf);
	num2 = int(buf);
	these_quotes.pop(buf);
	was_quote = True;
	goto(#these_quotes_s);
	#these_quotes_e:
	[print(""), ((op_num > num1)AND(op_num < num2)), goto(#is_op_end)];
	to_add = False;
	goto(#push_op_end);
	#is_op_end:	
	goto(#quotes_s);
	#quotes_e:
	[goto(#push_op_end), ((was_quote)OR(NOT(to_add))), print("")];
	res.push(op_num);
	was_quote = False;
	to_add = True;
	#push_op_end:
	print("");
	op_nums.pop(buf);
	goto(#op_nums_s);
	#op_nums_e:
	return res;
};

void replace_print(string command){
	stack s;
	string buf;

	s = ops(command, "print");

	s.pop(buf);
	#replace_s:
	[goto(#replace_e), ("end" == buf), print("")];
	println(command);
	println(buf);
	s.pop(buf);	
	goto(#replace_s);
	#replace_e:	

	print("");
};

void main(){
	string command;
	init();
	
	next_command(command);
	#main_s:
	[goto(#main_e), ("end" == command), print("")];
	replace_print(command);
	send_command(command);
	next_command(command);
	goto(#main_s);
	#main_e:
	finish();
};

main();
