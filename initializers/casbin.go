package initializers

// var ENFORCER *casbin.Enforcer

func ValidateAAA(role string) bool {
	if role == "admin" {
		return true
	} else if role == "reseller" {
		return true
	} else {
		return false
	}
}

// func LoadAAA(db *gorm.DB) {
// 	gormAdapter, err := gormadapter.NewAdapterByDB(db)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	ENFORCER, err = casbin.NewEnforcer("rbac_model.conf", gormAdapter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// dashboard

// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/dashboards/tlocation", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/dashboards/tlocation", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/dashboards/timeline", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/dashboards/timeline", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/dashboards", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/dashboards", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/orders/report/client/month", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/orders/report/client/month", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/orders/report/hotspot/month", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/orders/report/hotspot/month", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/payments/year/current", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/payments/year/current", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/payments/month/current", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/payments/month/current", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/payments/today/count", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/payments/today/count", "GET")
// 	}
// 	//add policy
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/users", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/users", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/register", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/register", "POST")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("reseller", "/api/v1/users", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("reseller", "/api/v1/users", "GET")
// 	}
// 	// PROFILES
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/profiles", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/profiles", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/profiles", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/profiles", "POST")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/profiles", "DELETE"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/profiles", "DELETE")
// 	}

// 	// PACKAGES
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/packages", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/packages", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/packages/detail", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/packages/detail", "GET")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/packages", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/packages", "POST")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/packages", "DELETE"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/packages", "DELETE")
// 	}

// 	// VOUCHERS
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/vouchers/batch", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/vouchers/batch", "POST")
// 	}

// 	// MEMBERS
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/members", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/members", "POST")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/members", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/members", "POST")
// 	}
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/members/subscribe", "POST"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/members/subscribe", "POST")
// 	}

// 	// PRINTER
// 	if hasPolicy := ENFORCER.HasPolicy("admin", "/api/v1/printer/print", "GET"); !hasPolicy {
// 		ENFORCER.AddPolicy("admin", "/api/v1/printer/print", "GET")
// 	}
// }
