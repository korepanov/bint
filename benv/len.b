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
	SET_DEST("benv/len_program.b");	
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

stack len_poses(string command){
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

	find = "len(";
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
			epos = func_end(command, (pos + 3));
			el.push(epos);
			el.push(pos);
			res.push(el);
			el = null;
			};
		};	
		s.pop(buf);	
	};

	return res;
};

string modify_command(string command, string sub_command, int bpos, int epos){
	string new_command;
	string buf;
	int command_len;
 
	new_command = command[0:bpos];
	new_command = (new_command + sub_command);
	command_len = len(command);
	buf = command[epos:command_len];
	new_command = (new_command + buf);
	
	return new_command;
};

void modify(){
	stack s;
	stack el;
	int bpos;
	int epos;
	string buf;
	int number;
	string snumber;
	string sub_command; 

	next_command(command);
	
	while (NOT("end" == command)){
		number = 0;
		s = len_poses(command);
		s.pop(el);
		el.pop(buf);

		while (NOT("end" == buf)){
			s = len_poses(command);
			s.pop(el);
			el.pop(buf);
			bpos = int(buf);
			el.pop(buf);
			epos = int(buf);
			epos = (epos + 1);
			snumber = str(number);
			buf = ("int $l" + snumber);
			send_command(buf);
			buf = command[bpos:epos];
			buf = ((("$l" + snumber) + "=") + buf);
			send_command(buf);
			sub_command = ("$l" + snumber);
			command = modify_command(command, sub_command, bpos, epos);
			number = (number + 1);
			s.pop(el);
			el.pop(buf);
		};
		send_command(command);

		for (int i; i = 0; i < number; i = (i + 1)){
			string b;
			snumber = str(i);
			b = (("UNDEFINE($l" + snumber) + ")");
			send_command(b);		
		};

		next_command(command);
	};
};

void main(){
	init();
	modify();
	finish(); 
	
	if ("benv/import_program.b" == root_source){
		copy("benv/len_program.b", "benv/import_program.b");
	}else{
		copy("benv/len_program.b", "benv/trace_program.b");
	};
	
	DEL_DEST("benv/len_program.b");
	string t;
	string t2;
	string s;
	s = str(0);
	t = (("benv/trace/trace_program" + s) + ".b");
	e = exists(t);
	for (int number; number = 1; e; number = (number + 1)){
		t = (("benv/trace/trace_program" + s) + ".b");
		SET_SOURCE(t);
		t = (("benv/trace/len_program" + s) + ".b");
		SET_DEST(t);
		modify();
		finish();
		t = (("benv/trace/len_program" + s) + ".b");
		t2 = (("benv/trace/trace_program" + s) + ".b");
		copy(t, t2);
		t = (("benv/trace/len_program" + s) + ".b"); 
		DEL_DEST(t);
		s = str(number);
		t = (("benv/trace/trace_program" + s) + ".b");
		e = exists(t);
	};
};

main();
