void println(string s){
	print(s);
	print("\n");
};

stack reverse(stack s){
	string buf;
	stack res; 
	
	s.pop(buf);
	#_reverse_s:
	[goto(#_reverse_e), ("end" == buf), print("")];
	res.push(buf);
	s.pop(buf);
	goto(#_reverse_s);
	#_reverse_e:
	return res;
};

stack indexes(string s, string sub_s){
	stack res;
	
	int i;
	int pointer; 
	int s_len;
	int s_len_old;
	int sub_len;
 	
	s_len = len(s);
	sub_len = len(sub_s);
	s_len_old = s_len;
	i = index(s, sub_s);
	pointer = i;
	#_indexes_s:
	[goto(#_indexes_e), (-1 == i), print("")];
	i = (i + (s_len_old - s_len));
	res.push(i);
	pointer = (pointer + sub_len);
	s = s[pointer:s_len];
	s_len = len(s);
	i = index(s, sub_s);
	pointer = i;
	goto(#_indexes_s);
	#_indexes_e:
	res = reverse(res);
	return res;
};
