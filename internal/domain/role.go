package domain

type Role struct {
	ID       uint   `gorm:"column:id_rol;primaryKey;autoIncrement"`
	RoleName string `gorm:"column:nombre;type:varchar(50);not null;uniqueIndex"`

	Users []User `gorm:"many2many:usuario_rol;foreignKey:ID;joinForeignKey:id_rol;References:ID;joinReferences:id_usuario"`
}

func (Role) TableName() string {
	return "rol"
}
