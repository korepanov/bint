string cut(string s, string sub_s){
	return (s[index(s, sub_s):len(s)]);
};


void main(){
	string s;
	s = "«Страна багровых туч» — приключенческая фантастическая повесть";
	s = (s + ", первое крупное произведение Аркадия и Бориса Стругацких. ");
	s = (s + "Замысел повести возник у Аркадия Стругацкого во время военной службы на Камчатке");
	s = (s + ", первые наброски и черновики датированы 1952—1954 годами. Основной текст написан");
	s = (s + " братьями Стругацкими в 1956—1957 годах и подвергся далее существенной редакционной правке.");
	s = (s + " Книга была подписана в печать в 1959 году, вызвала преимущественно положительную реакцию ");
	s = (s + "критиков и коллег-писателей. Аркадий и Борис Стругацкие сразу стали одними ");
	s = (s + "из самых известных советских фантастов.");

	print((cut(s, "1952") + "\n"));	 
	
};
main();
