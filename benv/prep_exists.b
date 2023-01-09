#import "stdlib/core.b"

string root_source;
string command;
bool e;

void init(){
	e = exists("benv/import_program.b");
	if (e){
		root_source = "benv/import_program.b";
	}else{
		root_source = "benv/trace_program.b";	
	};
	SET_SOURCE(root_source);
	SET_DEST("benv/prep_exists_program.b");	
};

void finish(){
	UNSET_SOURCE();
	UNSET_DEST();
};

void copy(string source, string dest){
	SET_SOURCE(source);
	SET_DEST(dest);
	string command;
	next_command(command); 
	while (NOT("end" == command)){
		send_command(command);
		next_command(command);	
	};

	UNSET_SOURCE();
	UNSET_DEST();
};

bool is_ex(string command){
	stack s;
	stack el;
	stack res;
	stack null; 
	int pos;
	int epos;
	string find;
	bool is_let;
	bool is_dig;
	string symbol;
	int buf_pos;

	find = "exists(";
	s = ops(command, find);
	string buf;
	s.pop(buf);
	
	while (NOT("end" == buf)){
		pos = int(buf);
		if (NOT(0 == pos)){
			buf_pos = (pos - 1);
			symbol = command[buf_pos]; 
			is_let = is_letter(symbol);
			is_dig = is_digit(symbol);
			
			if (NOT(((is_let)OR(is_dig))OR("_" == symbol))){
				return True;
			};
		};	
		s.pop(buf);	
	};

	return False;
};

void modify(){
	int bpos;
	int epos;
	int number;
	string snumber;
	string sub_command;
	int i; 
	
	next_command(command);
	
	while (NOT("end" == command)){
		number = 0;
		i = index(command, "{");
		if ((is_ex(command)) AND (NOT(-1 == i))){
			string t;
			int command_len;
			i = index(command, "{");
            i = (i + 1);
			t = command[0:i];
			send_command((t + "print(\"\")"));
			command_len = len(command);
			command = command[i:command_len]; 		
		};
		
		send_command(command); 
		next_command(command);
	};
};

void main(){
	init();
	modify();
	finish();
	
	if ("benv/import_program.b" == root_source){
		copy("benv/prep_exists_program.b", "benv/import_program.b");
	}else{
		copy("benv/prep_exists_program.b", "benv/trace_program.b");
	};
	
	DEL_DEST("benv/prep_exists_program.b");
	string t;
	string t2;
	string s;
	s = str(0);
	t = (("benv/trace/trace_program" + s) + ".b");
	e = exists(t);
	for (int number; number = 1; e; number = (number + 1)){
		t = (("benv/trace/trace_program" + s) + ".b");
		SET_SOURCE(t);
		t = (("benv/trace/prep_exists_program" + s) + ".b");
		SET_DEST(t);
		modify();
		finish();
		t = (("benv/trace/prep_exists_program" + s) + ".b");
		t2 = (("benv/trace/trace_program" + s) + ".b");
		copy(t, t2);
		t = (("benv/trace/prep_exists_program" + s) + ".b"); 
		DEL_DEST(t);
		s = str(number);
		t = (("benv/trace/trace_program" + s) + ".b");
		e = exists(t);
	};
};

main();
