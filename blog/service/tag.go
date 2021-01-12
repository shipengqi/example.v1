package service

func (s *Service) GetTags(maps map[string]interface{}) (data map[string]interface{}) {
	data = make(map[string]interface{})
	data["lists"], _ = s.dao.GetTags(0, 10, maps)
	data["total"] = s.dao.GetTagTotal(maps)
	return
}