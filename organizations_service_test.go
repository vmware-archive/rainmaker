package rainmaker_test

import (
	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OrganizationsService", func() {
	var (
		config       rainmaker.Config
		token        string
		service      rainmaker.OrganizationsService
		organization rainmaker.Organization
	)

	BeforeEach(func() {
		var err error

		token = "token"
		config = rainmaker.Config{
			Host: fakeCloudController.URL(),
		}
		client := rainmaker.NewClient(config)
		service = client.Organizations

		organization, err = service.Create("test-org", token)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		It("creates a new organization that can be fetched from the API", func() {
			organization, err := service.Create("my-new-org", token)
			Expect(err).NotTo(HaveOccurred())
			Expect(organization.Name).To(Equal("my-new-org"))

			fetchedOrg, err := service.Get(organization.GUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedOrg).To(Equal(organization))
		})
	})

	Describe("Get", func() {
		It("returns the organization matching the given GUID", func() {
			var err error
			orgGUID := organization.GUID

			organization, err = service.Get(orgGUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(organization.GUID).To(Equal(orgGUID))
		})
	})

	Context("when listing related users", func() {
		var user1, user2, user3 rainmaker.User

		BeforeEach(func() {
			var err error
			usersService := rainmaker.NewUsersService(config)

			user1, err = usersService.Create("user-123", token)
			Expect(err).NotTo(HaveOccurred())

			user2, err = usersService.Create("user-456", token)
			Expect(err).NotTo(HaveOccurred())

			user3, err = usersService.Create("user-789", token)
			Expect(err).NotTo(HaveOccurred())
		})

		Describe("ListUsers", func() {
			BeforeEach(func() {
				err := organization.Users.Associate(user1.GUID, token)
				Expect(err).NotTo(HaveOccurred())

				err = organization.Users.Associate(user2.GUID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the users belonging to the organization", func() {
				list, err := service.ListUsers(organization.GUID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(list.TotalResults).To(Equal(2))
				Expect(list.TotalPages).To(Equal(1))
				Expect(list.Users).To(HaveLen(2))

				var userGUIDs []string
				for _, user := range list.Users {
					userGUIDs = append(userGUIDs, user.GUID)
				}

				Expect(userGUIDs).To(ConsistOf([]string{user1.GUID, user2.GUID}))
			})

			Context("when the organization does not exist", func() {
				It("returns an error", func() {
					_, err := service.ListUsers("org-does-not-exist", token)
					Expect(err).To(BeAssignableToTypeOf(rainmaker.NotFoundError{}))
				})
			})
		})

		Describe("ListBillingManagers", func() {
			BeforeEach(func() {
				err := organization.BillingManagers.Associate(user2.GUID, token)
				Expect(err).NotTo(HaveOccurred())

				err = organization.BillingManagers.Associate(user3.GUID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the billing managers belonging to the organization", func() {
				list, err := service.ListBillingManagers(organization.GUID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(list.TotalResults).To(Equal(2))
				Expect(list.TotalPages).To(Equal(1))
				Expect(list.Users).To(HaveLen(2))

				var userGUIDs []string
				for _, user := range list.Users {
					userGUIDs = append(userGUIDs, user.GUID)
				}

				Expect(userGUIDs).To(ConsistOf([]string{user2.GUID, user3.GUID}))
			})
		})

		Describe("ListAuditors", func() {
			BeforeEach(func() {
				err := organization.Auditors.Associate(user1.GUID, token)
				Expect(err).NotTo(HaveOccurred())

				err = organization.Auditors.Associate(user3.GUID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the auditors belonging to the organization", func() {
				list, err := service.ListAuditors(organization.GUID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(list.TotalResults).To(Equal(2))
				Expect(list.TotalPages).To(Equal(1))
				Expect(list.Users).To(HaveLen(2))

				var userGUIDs []string
				for _, user := range list.Users {
					userGUIDs = append(userGUIDs, user.GUID)
				}

				Expect(userGUIDs).To(ConsistOf([]string{user1.GUID, user3.GUID}))
			})
		})

		Describe("ListManagers", func() {
			BeforeEach(func() {
				err := organization.Managers.Associate(user3.GUID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the managers belonging to the organization", func() {
				list, err := service.ListManagers(organization.GUID, token)
				Expect(err).NotTo(HaveOccurred())
				Expect(list.TotalResults).To(Equal(1))
				Expect(list.TotalPages).To(Equal(1))
				Expect(list.Users).To(HaveLen(1))

				Expect(list.Users[0].GUID).To(Equal(user3.GUID))
			})
		})

	})
})
