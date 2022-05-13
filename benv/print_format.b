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

	num1 = 0;
	num2 = len(command);
	num2 = (num2 - 1); 
	
	op_nums = indexes(command, op);	
	print(command);
	print("\n");

	quotes = reg_find("\"(.*?)\"", command);
	
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
	goto(#these_quotes_s);
	#these_quotes_e:
	print("\n");
	goto(#quotes_s);
	#quotes_e:

	return res;
};

void replace_print(string command){
	stack res;
	res = ops(command, "print");
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
